// SPDX-License-Identifier: GPL-3.0

package kubernetes

import (
	"errors"
	"fmt"

	pb "github.com/comrade-coop/apocryph/pkg/proto"
	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	policy "github.com/sigstore/policy-controller/pkg/apis/policy/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discovery "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8cl "sigs.k8s.io/controller-runtime/pkg/client"
)

func NewService(port *pb.Container_Port, portName string, httpSO *kedahttpv1alpha1.HTTPScaledObject, labels map[string]string) (*corev1.Service, int32, error) {

	servicePort := int32(port.ServicePort)
	if servicePort == 0 {
		servicePort = int32(port.ContainerPort)
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("tpod-svc-%v", port.Name),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Port:       servicePort,
				TargetPort: intstr.FromString(portName),
			}},
			Selector: labels,
		},
	}
	// since we don't use selectors, we need to manually create the endpoint
	// that routes to the tpod insights service

	switch ep := port.ExposedPort.(type) {
	case *pb.Container_Port_HostHttpHost:
		service.Spec.Type = corev1.ServiceTypeClusterIP
		service.Spec.Ports[0].Port = servicePort
		if len(httpSO.Spec.Hosts) > 0 {
			return nil, 0, errors.New("Multiple HTTP hosts in pod definition")
		}
		httpSO.Spec.Hosts = []string{ep.HostHttpHost}
	case *pb.Container_Port_HostTcpPort:
		service.Spec.Type = corev1.ServiceTypeNodePort
		service.Spec.Ports[0].Protocol = "TCP"
		if ep.HostTcpPort != 0 {
			service.Spec.Ports[0].NodePort = int32(ep.HostTcpPort)
		}
	}
	return service, servicePort, nil
}

func NewHttpSo(namespace, name string) *kedahttpv1alpha1.HTTPScaledObject {
	return &kedahttpv1alpha1.HTTPScaledObject{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: kedahttpv1alpha1.HTTPScaledObjectSpec{},
	}
}

func GetResource(kind string) k8cl.Object {
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
	case "ClusterImagePolicy":
		return &policy.ClusterImagePolicy{}

	case "EndpointSlice":
		return &discovery.EndpointSlice{}
	case "Endpoints":
		return &corev1.Endpoints{}

	}
	return nil
}

type AttestationValue struct {
	URL       string `json:"url"`
	Issuer    string `json:"issuer"`
	Identity  string `json:"identity"`
	Signature string `json:"signature"`
}
