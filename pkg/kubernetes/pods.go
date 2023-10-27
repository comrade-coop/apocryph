package kubernetes

import (
	"context"
	"fmt"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FetchSecret func(cid []byte) (map[string][]byte, error)

func ApplyPodRequest(ctx context.Context, client client.Client, podManifest *pb.Pod, response *pb.ProvisionPodResponse) error {
	labels := map[string]string{"tpod": "1"}

	startupReplicas := int32(0)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "tpod-dep-",
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

	httpSO := CreateHttpSo()

	localhostAliases := corev1.HostAlias{IP: "127.0.0.1"}

	for cIdx, container := range podManifest.Containers {
		containerSpec := corev1.Container{
			Name:       container.Name,
			Image:      container.Image.Url,
			Command:    container.Entrypoint,
			Args:       container.Command,
			WorkingDir: container.WorkingDir,
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
			service, servicePort, err := GetService(port, portName, httpSO, labels)
			if err != nil {
				return err
			}
			if err := client.Create(ctx, service); err != nil {
				return err
			}

			multiaddrPart := ""

			switch port.ExposedPort.(type) {
			case *pb.Container_Port_HostHttpHost:
				httpSO.Spec.ScaleTargetRef.Service = service.ObjectMeta.Name
				httpSO.Spec.ScaleTargetRef.Port = servicePort
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
	}
	podTemplate.Spec.HostAliases = append(podTemplate.Spec.HostAliases, localhostAliases)
	for _, volume := range podManifest.Volumes {
		volumeSpec := corev1.Volume{
			Name: volume.Name,
		}
		switch volume.Type {
		case pb.Volume_VOLUME_EMPTY:
			volumeSpec.VolumeSource.EmptyDir = &corev1.EmptyDirVolumeSource{}
		case pb.Volume_VOLUME_FILESYSTEM:
			persistentVolumeClaim := &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "tpod-pvc-",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					Resources: corev1.ResourceRequirements{
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

			if err := client.Create(ctx, persistentVolumeClaim); err != nil {
				return err
			}

			volumeSpec.VolumeSource.PersistentVolumeClaim = &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: persistentVolumeClaim.ObjectMeta.Name,
			}
		case pb.Volume_VOLUME_SECRET:
			secretBytes := volume.GetSecret().Contents

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "tpod-secret-",
				},
				Data: map[string][]byte{
					"data": secretBytes,
				},
			}

			if err := client.Create(ctx, secret); err != nil {
				return err
			}

			volumeSpec.VolumeSource.Secret = &corev1.SecretVolumeSource{
				SecretName: secret.ObjectMeta.Name,
			}
		}
		podTemplate.Spec.Volumes = append(podTemplate.Spec.Volumes, volumeSpec)
	}

	if err := client.Create(ctx, deployment); err != nil {
		return err
	}

	if httpSO.Spec.ScaleTargetRef.Service != "" {
		httpSO.Spec.ScaleTargetRef.Deployment = deployment.ObjectMeta.Name
		minReplicas := int32(podManifest.Replicas.Min)
		maxReplicas := int32(podManifest.Replicas.Max)
		targetPendingRequests := int32(podManifest.Replicas.TargetPendingRequests)
		httpSO.Spec.Replicas = &kedahttpv1alpha1.ReplicaStruct{
			Min: &minReplicas,
		}
		if maxReplicas > 0 {
			httpSO.Spec.Replicas.Max = &maxReplicas
		}
		if targetPendingRequests > 0 {
			httpSO.Spec.TargetPendingRequests = &targetPendingRequests
		}

		if err := client.Create(ctx, httpSO); err != nil {
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
