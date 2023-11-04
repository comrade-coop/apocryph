package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8cl "sigs.k8s.io/controller-runtime/pkg/client"
)

type FetchSecret func(cid []byte) (map[string][]byte, error)

func GetRessource(kind string) interface{} {
	switch kind {
	case "Service":
		return &corev1.Service{}
	case "Volume":
		return &corev1.PersistentVolumeClaim{}
	case "Secret":
		return &corev1.Secret{}
	case "Deployment":
		return &appsv1.Deployment{}
	case "HttpSo":
		return &kedahttpv1alpha1.HTTPScaledObject{}
	}
	return nil
}

func patchOrCreate(ctx context.Context, ressourceName, kind, namespace string, ressource interface{}, client k8cl.Client, patch bool) error {
	if patch == true {
		key := &k8cl.ObjectKey{
			Namespace: namespace,
			Name:      ressourceName,
		}
		oldRessource := GetRessource(kind)

		err := client.Get(ctx, *key, oldRessource.(k8cl.Object))
		if err != nil {
			return err
		}

		patch, err := json.Marshal(ressource)
		if err != nil {
			return err
		}

		err = client.Patch(ctx, oldRessource.(k8cl.Object), k8cl.RawPatch(types.MergePatchType, patch))
		if err != nil {
			return err
		}
		return nil
	}
	if err := client.Create(ctx, ressource.(k8cl.Object)); err != nil {
		return err
	}
	return nil
}

func cleanNamespace(ctx context.Context, namespace string, activeRessources []string, client k8cl.Client) error {
	kindList := []string{"Service", "Volume", "Secret", "Deployment", "HttpSo"}
	fmt.Printf("Active Ressources: %v \n", activeRessources)
	for _, kind := range kindList {
		switch kind {
		case "Service":
			list := &corev1.ServiceList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeRessources, rsrc.GetName()) && strings.Contains(rsrc.GetName(), "tpod") {
					fmt.Printf("Deleting Service %v:%v \n", i, rsrc.GetName())
					client.Delete(ctx, &rsrc)
				}
			}
		case "Volume":
			list := &corev1.PersistentVolumeList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				// fmt.Printf("volume: %v\n", rsrc.GetName()) // where do all these come from and where is tpod-pvc-db-data?
				if !slices.Contains(activeRessources, rsrc.GetName()) && strings.Contains(rsrc.GetName(), "tpod") {
					fmt.Printf("Deleting Volume %v: %v \n", i, rsrc.GetName())
					client.Delete(ctx, &rsrc)
				}
			}
		case "Secret":
			list := &corev1.SecretList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeRessources, rsrc.GetName()) && strings.Contains(rsrc.GetName(), "tpod") {
					fmt.Printf("Deleting Secret %v: %v \n", i, rsrc.GetName())
					client.Delete(ctx, &rsrc)
				}
			}
		case "Deployment":
			list := &appsv1.DeploymentList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeRessources, rsrc.GetName()) && strings.Contains(rsrc.GetName(), "tpod") {
					fmt.Printf("Deleting Deployment %v: %v \n", i, rsrc.GetName())
					client.Delete(ctx, &rsrc)
				}
			}
		case "HttpSo":
			list := &kedahttpv1alpha1.HTTPScaledObjectList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeRessources, rsrc.GetName()) && strings.Contains(rsrc.GetName(), "tpod") {
					fmt.Printf("Deleting HttpSo %v: %v \n", i, rsrc.GetName())
					client.Delete(ctx, &rsrc)
				}
			}
		}

	}
	return nil

}

func ApplyPodRequest(ctx context.Context, client k8cl.Client, podManifest *pb.Pod, patch bool, namespace string, response *pb.ProvisionPodResponse) error {
	labels := map[string]string{"tpod": "1"}
	activeRessource := []string{}
	startupReplicas := int32(0)
	var deploymentName = fmt.Sprintf("tpod-dep-%v", namespace)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
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

	httpSoName := fmt.Sprintf("tpod-httpso-%v", namespace)
	httpSO := NewHttpSo(namespace, httpSoName)

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
			service, servicePort, err := NewService(port, portName, httpSO, labels)
			if err != nil {
				return err
			}

			err = patchOrCreate(ctx, service.GetName(), "Service", namespace, service, client, patch)
			if err != nil {
				return err
			}

			activeRessource = append(activeRessource, service.GetName())
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

			err := patchOrCreate(ctx, volumeName, "Volume", namespace, persistentVolumeClaim, client, patch)
			if err != nil {
				return err
			}
			activeRessource = append(activeRessource, volumeName)

			volumeSpec.VolumeSource.PersistentVolumeClaim = &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: persistentVolumeClaim.ObjectMeta.Name,
			}
		case pb.Volume_VOLUME_SECRET:
			var secretName = fmt.Sprintf("tpod-secret-%v", volume.Name)
			secretBytes := volume.GetSecret().Contents

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: secretName,
				},
				Data: map[string][]byte{
					"data": secretBytes,
				},
			}

			err := patchOrCreate(ctx, secretName, "Secret", namespace, secret, client, patch)
			if err != nil {
				return err
			}
			activeRessource = append(activeRessource, secretName)

			volumeSpec.VolumeSource.Secret = &corev1.SecretVolumeSource{
				SecretName: secret.ObjectMeta.Name,
			}
		}
		podTemplate.Spec.Volumes = append(podTemplate.Spec.Volumes, volumeSpec)
	}

	err := patchOrCreate(ctx, deploymentName, "Deployment", namespace, deployment, client, patch)
	if err != nil {
		return err
	}
	activeRessource = append(activeRessource, deploymentName)

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

		err := patchOrCreate(ctx, httpSoName, "HttpSo", namespace, httpSO, client, patch)
		if err != nil {
			return err
		}
		activeRessource = append(activeRessource, httpSoName)

	}
	if patch == true {
		err := cleanNamespace(ctx, namespace, activeRessource, client)
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
