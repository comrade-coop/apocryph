// SPDX-License-Identifier: GPL-3.0

package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

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

var proxyImageReference = "ttl.sh/comradecoop/apocryph/tpod-proxy@sha256:fc671075e920aa8022e65b330cd8e2b56b49f9b46e4b96a958238c63c65a84f2"

// TODO: use actual comrade account details
var proxyVerificationDetails = &pb.VerificationDetails{
	Identity: "comradecoop@email.com",
	Issuer:   "https://github.com/login/oauth",
}

func CreateTpodProxyPolicy(ctx context.Context, client k8cl.Client) error {
	policyName := fmt.Sprintf("tpod-policy-proxy")
	sigstorePolicy := &policy.ClusterImagePolicy{
		TypeMeta: metav1.TypeMeta{Kind: "ClusterImagePolicy"}, ObjectMeta: metav1.ObjectMeta{Name: policyName},
		Spec: policy.ClusterImagePolicySpec{
			Images:      []policy.ImagePattern{{Glob: proxyImageReference}},
			Authorities: []policy.Authority{{Keyless: &policy.KeylessRef{}}},
		}}
	identity := policy.Identity{Issuer: proxyVerificationDetails.Issuer, Subject: proxyVerificationDetails.Identity}
	sigstorePolicy.Spec.Authorities[0].Keyless.Identities = []policy.Identity{identity}
	err := updateOrCreate(ctx, policyName, "ClusterImagePolicy", "default", sigstorePolicy, client, false)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		log.Println("warning: Policies were not deleted properly")
	} else {
		return err
	}
	return nil
}

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
