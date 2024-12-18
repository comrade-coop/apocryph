// SPDX-License-Identifier: GPL-3.0

package kubernetes

import (
	"context"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/comrade-coop/apocryph/pkg/constants"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	policy "github.com/sigstore/policy-controller/pkg/apis/policy/v1beta1"
	"golang.org/x/crypto/sha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8cl "sigs.k8s.io/controller-runtime/pkg/client"
)

type FetchSecret func(cid []byte) (map[string][]byte, error)

// given a k8s resource; it checks the existence of that resource in the
// cluster, if it exists it Will update it if needed, if not it will create it
func updateOrCreate(ctx context.Context, resourceName, kind, namespace string, resource k8cl.Object, client k8cl.Client, update bool) error {
	if update {
		key := &k8cl.ObjectKey{
			Namespace: namespace,
			Name:      resourceName,
		}
		oldResource := GetResource(kind)

		resource.SetNamespace(namespace)
		resource.SetName(resourceName)

		err := client.Get(ctx, *key, oldResource)
		resource.SetResourceVersion(oldResource.GetResourceVersion()) // resource version should be retrieved from the old resource in order for httpSo to work
		if err != nil {
			log.Printf("Added New Resource: %v \n", resourceName)
			if err := client.Create(ctx, resource); err != nil {
				return fmt.Errorf("Failed creating resource:%v: %v\n", resourceName, err)
			}
			return nil
		}

		err = client.Update(ctx, resource)
		if err != nil {
			return fmt.Errorf("Failed updating resource:%v, %v\n", resourceName, err)
		}
		log.Printf("Updated %v \n", resourceName)
		return nil
	}
	if err := client.Create(ctx, resource.(k8cl.Object)); err != nil {
		return fmt.Errorf("Failed creating resource %v:%v\n", resourceName, err)
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
	proxyImage string,
) error {
	if podManifest == nil {
		return fmt.Errorf("Expected value for pod")
	}

	labels := map[string]string{"tpod": "1"}
	podLabels := map[string]string{}
	for k, v := range labels {
		podLabels[k] = v
	}
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
					Labels: podLabels,
				},
			},
		},
	}
	attestationValues := []AttestationValue{}

	podTemplate := &deployment.Spec.Template

	httpSoName := fmt.Sprintf("tpod-hso")
	httpSO := NewHttpSo(namespace, httpSoName)

	localhostAliases := corev1.HostAlias{IP: "127.0.0.1"}

	var httpPortName string

	for cIdx, container := range podManifest.Containers {
		containerSpec := corev1.Container{
			Name:            container.Name,
			Image:           images[container.Name],
			ImagePullPolicy: corev1.PullIfNotPresent, // make sure images are pulled localy (didn't never pull for the tpod-proxy image, for now)
			Command:         container.Entrypoint,
			Args:            container.Command,
			WorkingDir:      container.WorkingDir,
		}
		if container.GetImage().GetCid() == nil {
			containerSpec.ImagePullPolicy = corev1.PullIfNotPresent // HACK: Allow URL-only images for easier testing
		}

		if podManifest.KeyPair != nil {
			// save as hex to parse later as hex
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PAYMENT_ADDR_KEY, Value: common.BytesToAddress(paymentChannel.ContractAddress).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PUBLISHER_ADDR_KEY, Value: common.BytesToAddress(paymentChannel.PublisherAddress).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PROVIDER_ADDR_KEY, Value: common.BytesToAddress(paymentChannel.ProviderAddress).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.POD_ID_KEY, Value: common.BytesToHash(paymentChannel.PodID).Hex()})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PUBLIC_ADDRESS_KEY, Value: podManifest.KeyPair.PubAddress})
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: constants.PRIVATE_KEY, Value: podManifest.KeyPair.PrivateKey})
		}

		attestationValues = append(attestationValues, AttestationValue{
			URL:       container.GetImage().GetUrl(),
			Signature: container.GetImage().VerificationDetails.GetSignature(),
			Issuer:    container.GetImage().VerificationDetails.GetIssuer(),
			Identity:  container.GetImage().VerificationDetails.GetIdentity(),
		})

		if podManifest.VerificationSettings.GetForcePolicy() && container.Image.VerificationDetails != nil {
			policyName := fmt.Sprintf("policy-%v-%v", podId, cIdx) // NOTE: Will break for two-digit container ids (see note around podId)
			sigstorePolicy := &policy.ClusterImagePolicy{
				TypeMeta: metav1.TypeMeta{Kind: "ClusterImagePolicy"},
				ObjectMeta: metav1.ObjectMeta{
					Name: policyName,
					Labels: map[string]string{
						LabelClusterImagePolicy: podId,
					},
				},
				Spec: policy.ClusterImagePolicySpec{
					Images:      []policy.ImagePattern{{Glob: container.Image.Url}},
					Authorities: []policy.Authority{{Keyless: &policy.KeylessRef{}}},
				},
			}
			identity := policy.Identity{Issuer: container.Image.VerificationDetails.Issuer, Subject: container.Image.VerificationDetails.Identity}
			sigstorePolicy.Spec.Authorities[0].Keyless.Identities = []policy.Identity{identity}
			err := updateOrCreate(ctx, policyName, "ClusterImagePolicy", namespace, sigstorePolicy, client, update)
			if err != nil && strings.Contains(err.Error(), "already exists") {
				log.Println("warning: Policies were not deleted properly")
				continue
			} else {
				return err
			}
		}

		for field, value := range container.Env {
			containerSpec.Env = append(containerSpec.Env, corev1.EnvVar{Name: field, Value: value})
		}

		for _, port := range container.Ports {
			if int32(port.ContainerPort) == constants.ReservedTpodProxyPort {
				return fmt.Errorf("Port %d is reserved", port.ContainerPort)
			}
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
				httpPortName = portName
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

	// Public verifiability configured; deploy the proxy
	if podManifest.VerificationSettings.GetPublicVerifiability() == true {
		if httpPortName == "" {
			return fmt.Errorf("Public verifiability requires having an HTTP-exposed port")
		}

		// TODO: Do not delete the previous versionedService(s) until the new deployment is fully rolled out
		versionedServiceName := getHttpBackingServiceName(attestationValues)

		versionedLabels := map[string]string{"tpod-version": versionedServiceName}
		for k, v := range versionedLabels {
			podLabels[k] = v
		}

		versionedService := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: versionedServiceName,
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Port:       80,
					TargetPort: intstr.FromString(httpPortName),
				}},
				Selector: versionedLabels,
			},
		}
		err := updateOrCreate(ctx, versionedServiceName, "Service", namespace, versionedService, client, update)
		if err != nil {
			return err
		}
		activeResource = append(activeResource, versionedServiceName)
		attestationValuesJson, err := json.Marshal(attestationValues)
		if err != nil {
			return fmt.Errorf("Failed to marshal attestation values: %v", err)
		}

		proxyContainer := corev1.Container{
			Name:  "tpodsproxy-sc",
			Ports: []corev1.ContainerPort{{ContainerPort: constants.ReservedTpodProxyPort, Name: "tpodsproxy-sc"}},
			Image: proxyImage,
			Args: []string{
				versionedServiceName,
				"." + namespace + ".svc.cluster.local",
				string(attestationValuesJson),
				"--address",
				fmt.Sprintf(":%d", constants.ReservedTpodProxyPort),
			},
		}
		podTemplate.Spec.Containers = append(podTemplate.Spec.Containers, proxyContainer)

		serviceProxyPort := constants.ReservedTpodProxyPort
		serviceProxyName := "tpod-svc-proxy"
		serviceProxy := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceProxyName,
			},
			Spec: corev1.ServiceSpec{
				Type:     corev1.ServiceTypeClusterIP,
				Ports:    []corev1.ServicePort{{Port: serviceProxyPort, TargetPort: intstr.FromInt32(constants.ReservedTpodProxyPort)}},
				Selector: labels,
			},
		}
		err = updateOrCreate(ctx, serviceProxyName, "Service", namespace, serviceProxy, client, update)
		if err != nil {
			return err
		}
		activeResource = append(activeResource, serviceProxyName)
		httpSO.Spec.ScaleTargetRef.Service = serviceProxyName
		httpSO.Spec.ScaleTargetRef.Port = serviceProxyPort
	}

	if httpSO.Spec.ScaleTargetRef.Service == "" { // No scaler configured - just deploy min replicas
		startupReplicas = int32(podManifest.Replicas.GetMin())
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

func getHttpBackingServiceName(attestationValues []AttestationValue) string {
	hash := sha3.New256()
	for _, v := range attestationValues {
		hash.Write([]byte(v.URL))
		hash.Write([]byte{0})
	}
	resultHash := hash.Sum(nil)
	return "tpod-vh-" + strings.TrimRight(strings.ToLower(base32.StdEncoding.EncodeToString(resultHash)), "=")
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
		case *pb.Resource_AmountKibi:
			quantity = *resource.NewQuantity(int64(q.AmountKibi)*1024, resource.BinarySI)
		case *pb.Resource_AmountMebi:
			quantity = *resource.NewQuantity(int64(q.AmountMebi)*1024*1024, resource.BinarySI)
		case *pb.Resource_AmountGibi:
			quantity = *resource.NewQuantity(int64(q.AmountGibi)*1024*1024*1024, resource.BinarySI)
		}
		result[corev1.ResourceName(res.Resource)] = quantity
	}
	return result
}
