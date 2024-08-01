// SPDX-License-Identifier: GPL-3.0

package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/comrade-coop/apocryph/pkg/constants"
	tpcrypto "github.com/comrade-coop/apocryph/pkg/crypto"
	"github.com/comrade-coop/apocryph/pkg/ethereum"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8cl "sigs.k8s.io/controller-runtime/pkg/client"
)

type FetchSecret func(cid []byte) (map[string][]byte, error)

const PRIVATE_KEY = "PRIVATE_KEY"
const PUBLIC_ADDRESS = "PUBLIC_ADDRESS"

func updateOrCreate(ctx context.Context, resourceName, kind, namespace string, resource interface{}, client k8cl.Client, update bool) error {
	if update {
		key := &k8cl.ObjectKey{
			Namespace: namespace,
			Name:      resourceName,
		}
		oldResource := GetResource(kind)

		updatedResource := resource.(k8cl.Object)
		updatedResource.SetNamespace(namespace)
		updatedResource.SetName(resourceName)

		err := client.Get(ctx, *key, oldResource.(k8cl.Object))
		updatedResource.SetResourceVersion(oldResource.(k8cl.Object).GetResourceVersion()) // resource version should be retrieved from the old resource in order for httpSo to work
		if err != nil {
			fmt.Printf("Added New Resource: %v \n", resourceName)
			if err := client.Create(ctx, updatedResource); err != nil {
				return err
			}
			return nil
		}

		err = client.Update(ctx, updatedResource)
		if err != nil {
			return err
		}
		fmt.Printf("Updated %v \n", resourceName)
		return nil
	}
	if err := client.Create(ctx, resource.(k8cl.Object)); err != nil {
		return err
	}
	return nil
}

func ApplyPodRequest(
	ctx context.Context,
	client k8cl.Client,
	namespace string,
	update bool,
	podManifest *pb.Pod,
	paymentChannel *pb.PaymentChannel,
	images map[string]string,
	secrets map[string][]byte,
	response *pb.ProvisionPodResponse,
) error {
	if podManifest == nil {
		return fmt.Errorf("Expected value for pod")
	}

	labels := map[string]string{"tpod": "1"}
	depLabels := map[string]string{}
	activeResource := []string{}
	startupReplicas := int32(0)
	podId := strings.Split(namespace, "-")[1]
	// NOTE: podId is 52 characters long. Most names are 63 characters max. Thus: don't make any constant string longer than 10 characters. "tpod-dep-" is 9.
	var deploymentName = fmt.Sprintf("tpod-dep-%v", podId)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   deploymentName,
			Labels: depLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &startupReplicas,
			Selector: metav1.SetAsLabelSelector(labels),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
			},
		},
	}

	podTemplate := &deployment.Spec.Template

	httpSoName := fmt.Sprintf("tpod-hso")
	httpSO := NewHttpSo(namespace, httpSoName)

	localhostAliases := corev1.HostAlias{IP: "127.0.0.1"}

	for cIdx, container := range podManifest.Containers {
		containerSpec := corev1.Container{
			Name:            container.Name,
			Image:           images[container.Name],
			ImagePullPolicy: corev1.PullNever, // make sure images are pulled localy
			Command:         container.Entrypoint,
			Args:            container.Command,
			WorkingDir:      container.WorkingDir,
		}

		if podManifest.KeyPair != nil {
			privatekey, err := tpcrypto.DecryptWithKey(podManifest.KeyPair.Key, podManifest.KeyPair.PrivateKey)
			if err != nil {
				return fmt.Errorf("Failed Decrypting private key: %v\n", err)
			}

			key, err := ethereum.EncodePrivateKey(privatekey)
			if err != nil {
				return fmt.Errorf("Failed encoding private key: %v\n", err)
			}

			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PRIVATE_KEY, Value: key})
			// save as hex to parse later as hex
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PAYMENT_ADDR_KEY, Value: common.BytesToAddress(paymentChannel.ContractAddress).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PUBLISHER_ADDR_KEY, Value: common.BytesToAddress(paymentChannel.PublisherAddress).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PROVIDER_ADDR_KEY, Value: common.BytesToAddress(paymentChannel.ProviderAddress).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.POD_ID_KEY, Value: common.BytesToHash(paymentChannel.PodID).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PRIVATE_KEY, Value: key})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PUBLIC_ADDRESS_KEY, Value: podManifest.KeyPair.PubAddress})
		}

		for field, value := range container.Env {
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: field, Value: value})
		}

		for _, port := range container.Ports {
			portName := fmt.Sprintf("p%d-%d", cIdx, port.ContainerPort)
			containerSpec.Ports = append(containerSpec.Ports, corev1.ContainerPort{
				ContainerPort: int32(port.ContainerPort),
				Name:          portName,
			})
			service, servicePort, err := NewService(port, portName, httpSO, labels)
			if err != nil {
				return err
			}

			err = updateOrCreate(ctx, service.GetName(), "Service", namespace, service, client, update)
			if err != nil {
				return err
			}

			activeResource = append(activeResource, service.GetName())
			multiaddrPart := ""

			switch port.ExposedPort.(type) {
			case *pb.Container_Port_HostHttpHost:
				httpSO.Spec.ScaleTargetRef.Service = service.ObjectMeta.Name
				httpSO.Spec.ScaleTargetRef.Port = servicePort
				httpSO.Spec.ScaleTargetRef.APIVersion = "apps/v1"

				multiaddrPart = fmt.Sprintf("http/%s", httpSO.Spec.Hosts[0])
			case *pb.Container_Port_HostTcpPort:
				multiaddrPart = fmt.Sprintf("tcp/%d", service.Spec.Ports[0].NodePort)
			}
			response.Addresses = append(response.Addresses, &pb.ProvisionPodResponse_ExposedHostPort{
				Multiaddr:     multiaddrPart,
				ContainerName: container.Name,
				ContainerPort: port.ContainerPort,
			})
		}
		for _, volume := range container.Volumes {
			volumeMount := corev1.VolumeMount{
				Name:      volume.Name,
				MountPath: volume.MountPath,
			}
			for _, targetVolume := range podManifest.Volumes {
				if targetVolume.Name == volume.Name {
					if targetVolume.Type == pb.Volume_VOLUME_SECRET {
						volumeMount.SubPath = "data" // NOTE: Change when secrets start supporting filesystems
					}
				}
			}
			containerSpec.VolumeMounts = append(containerSpec.VolumeMounts, volumeMount)
		}
		containerSpec.Resources.Requests = convertResourceList(container.ResourceRequests)
		// TODO: Enforce specifying resources?
		podTemplate.Spec.Containers = append(podTemplate.Spec.Containers, containerSpec)
		localhostAliases.Hostnames = append(localhostAliases.Hostnames, container.Name)
		if depLabels["containers"] == "" {
			depLabels["containers"] = containerSpec.Name
		} else {
			depLabels["containers"] = depLabels["containers"] + "_" + containerSpec.Name
		}
	}
	podTemplate.Spec.HostAliases = append(podTemplate.Spec.HostAliases, localhostAliases)
	for _, volume := range podManifest.Volumes {
		volumeSpec := corev1.Volume{
			Name: volume.Name,
		}
		var volumeName = fmt.Sprintf("tpod-pvc-%v", volume.Name)
		switch volume.Type {
		case pb.Volume_VOLUME_EMPTY:
			volumeSpec.VolumeSource.EmptyDir = &corev1.EmptyDirVolumeSource{}
		case pb.Volume_VOLUME_FILESYSTEM:
			persistentVolumeClaim := &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: volumeName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					Resources: corev1.VolumeResourceRequirements{
						Requests: convertResourceList(volume.GetFilesystem().ResourceRequests),
					},
				},
			}

			switch volume.AccessMode {
			case pb.Volume_VOLUME_RW_ONE:
				persistentVolumeClaim.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
			case pb.Volume_VOLUME_RW_MANY:
				persistentVolumeClaim.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany}
			}

			// pvcs couldn't and shouldn't be updated
			if !update {
				err := updateOrCreate(ctx, volumeName, "Volume", namespace, persistentVolumeClaim, client, update)
				if err != nil {
					return err
				}

			}
			activeResource = append(activeResource, volumeName)

			volumeSpec.VolumeSource.PersistentVolumeClaim = &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: persistentVolumeClaim.ObjectMeta.Name,
			}
		case pb.Volume_VOLUME_SECRET:
			var secretName = fmt.Sprintf("tpod-secret-%v", volume.Name)

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: secretName,
				},
				Data: map[string][]byte{
					"data": secrets[volume.Name],
				},
			}

			err := updateOrCreate(ctx, secretName, "Secret", namespace, secret, client, update)
			if err != nil {
				return err
			}
			activeResource = append(activeResource, secretName)

			volumeSpec.VolumeSource.Secret = &corev1.SecretVolumeSource{
				SecretName: secret.ObjectMeta.Name,
			}
		}
		podTemplate.Spec.Volumes = append(podTemplate.Spec.Volumes, volumeSpec)
	}
	err := updateOrCreate(ctx, deploymentName, "Deployment", namespace, deployment, client, update)
	if err != nil {
		return err
	}
	activeResource = append(activeResource, deploymentName)

	if httpSO.Spec.ScaleTargetRef.Service != "" {
		httpSO.Spec.ScaleTargetRef.Kind = "Deployment"
		httpSO.Spec.ScaleTargetRef.Name = deploymentName
		minReplicas := int32(podManifest.Replicas.GetMin())
		maxReplicas := int32(podManifest.Replicas.GetMax())
		targetPendingRequests := int32(podManifest.Replicas.GetTargetPendingRequests())
		httpSO.Spec.Replicas = &kedahttpv1alpha1.ReplicaStruct{
			Min: &minReplicas,
		}
		if maxReplicas > 0 {
			httpSO.Spec.Replicas.Max = &maxReplicas
		}
		if targetPendingRequests > 0 {
			httpSO.Spec.TargetPendingRequests = &targetPendingRequests
		}
		// update status
		httpSO.Status = kedahttpv1alpha1.HTTPScaledObjectStatus{}
		httpSO.Status.TargetWorkload = fmt.Sprintf("%s/%s/%s", httpSO.Spec.ScaleTargetRef.APIVersion, httpSO.Spec.ScaleTargetRef.Kind, httpSO.Spec.ScaleTargetRef.Name)
		httpSO.Status.TargetService = fmt.Sprintf("%s:%d", httpSO.Spec.ScaleTargetRef.Service, httpSO.Spec.ScaleTargetRef.Port)

		err := updateOrCreate(ctx, httpSoName, "HttpSo", namespace, httpSO, client, update)
		if err != nil {
			return err
		}
		activeResource = append(activeResource, httpSoName)

	}
	if update == true {
		err := cleanNamespace(ctx, namespace, activeResource, client)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertResourceList(resources []*pb.Resource) corev1.ResourceList {
	result := make(corev1.ResourceList, len(resources))
	for _, res := range resources {
		var quantity resource.Quantity
		switch q := res.Quantity.(type) {
		case *pb.Resource_Amount:
			quantity = *resource.NewQuantity(int64(q.Amount), resource.BinarySI)
		case *pb.Resource_AmountMillis:
			quantity = *resource.NewMilliQuantity(int64(q.AmountMillis), resource.BinarySI)
		}
		result[corev1.ResourceName(res.Resource)] = quantity
	}
	return result
}
