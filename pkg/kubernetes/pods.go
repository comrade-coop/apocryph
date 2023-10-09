package kubernetes

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ApplyPodRequest(ctx context.Context, client client.Client, manifest *pb.ProvisionPodRequest, response *pb.ProvisionPodResponse) error {
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

	httpSO := &kedahttpv1alpha1.HTTPScaledObject{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "tpod-svc",
		},
		Spec: kedahttpv1alpha1.HTTPScaledObjectSpec{},
	}

	for cIdx, container := range manifest.PodManifest.Containers {
		containerSpec := corev1.Container{
			Name:       container.Name,
			Image:      "nginxdemos/nginx-hello", // TODO
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
			servicePort := int32(port.ServicePort)
			if servicePort == 0 {
				servicePort = int32(port.ContainerPort)
			}

			service := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "tpod-svc-",
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{{
						Port:       servicePort,
						TargetPort: intstr.FromString(portName),
					}},
					Selector: labels,
				},
			}

			switch ep := port.ExposedPort.(type) {
			case *pb.Container_Port_HostHttpHost:
				service.Spec.Type = corev1.ServiceTypeClusterIP
				if len(httpSO.Spec.Hosts) > 0 {
					return errors.New("Multiple HTTP hosts in pod definition")
				}
				httpSO.Spec.Hosts = []string{ep.HostHttpHost}
			case *pb.Container_Port_HostTcpPort:
				service.Spec.Type = corev1.ServiceTypeNodePort
				service.Spec.Ports[0].Protocol = "TCP"
				if ep.HostTcpPort != 0 {
					service.Spec.Ports[0].NodePort = int32(ep.HostTcpPort)
				}
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
			containerSpec.VolumeMounts = append(containerSpec.VolumeMounts, corev1.VolumeMount{
				Name:      fmt.Sprintf("vol-%d", volume.VolumeIdx),
				MountPath: volume.MountPath,
			})
		}
		containerSpec.Resources.Requests = convertResourceList(container.ResourceRequests)
		// TODO: Enforce specifying resources?
		podTemplate.Spec.Containers = append(podTemplate.Spec.Containers, containerSpec)
	}
	for idx, volume := range manifest.PodManifest.Volumes {
		volumeSpec := corev1.Volume{
			Name: fmt.Sprintf("vol-%d", idx),
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
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "tpod-secret-",
				},
				Data: map[string][]byte{
					"cid": []byte(volume.GetSecret().Cid), // TODO
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
		minReplicas := int32(manifest.PodManifest.Replicas.Min)
		maxReplicas := int32(manifest.PodManifest.Replicas.Max)
		targetPendingRequests := int32(manifest.PodManifest.Replicas.TargetPendingRequests)
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
