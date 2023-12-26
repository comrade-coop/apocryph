// SPDX-License-Identifier: GPL-3.0

package kubernetes

import (
	"context"
	"fmt"

	pb "github.com/comrade-coop/apocryph/pkg/proto"
	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/encoding/protojson"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8cl "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LabelTrustedPodsNamespace            string = "coop.comrade/apocryph-namespace"
	AnnotationsTrustedPodsPaymentChannel string = "coop.comrade/apocryph-payment-contract"
	LabelIpfsP2P                         string = "coop.comrade/apocryph-p2p-helper"
	AnnotationsIpfsP2P                   string = "coop.comrade/apocryph-p2p-helper"
)

var TrustedPodsNamespaceFilter = client.HasLabels{LabelTrustedPodsNamespace}

func NewTrustedPodsNamespace(name string, paymentChannel *pb.PaymentChannel) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				LabelTrustedPodsNamespace: "true",
			},
		},
	}
	if paymentChannel != nil {
		namespace.ObjectMeta.Annotations = map[string]string{
			AnnotationsTrustedPodsPaymentChannel: protojson.Format(paymentChannel),
		}
	}
	return namespace
}

func TrustedPodsNamespaceGetChannel(namespace *corev1.Namespace) (*pb.PaymentChannel, error) {
	paymentChannel := &pb.PaymentChannel{}
	paymentChannelJson, ok := namespace.ObjectMeta.Annotations[AnnotationsTrustedPodsPaymentChannel]
	if !ok {
		return nil, nil
	}
	err := protojson.Unmarshal([]byte(paymentChannelJson), paymentChannel)
	if err != nil {
		return nil, err
	}
	return paymentChannel, nil
}

func cleanNamespace(ctx context.Context, namespace string, activeResources []string, client k8cl.Client) error {
	kindList := []string{"Service", "Volume", "Secret", "Deployment", "HttpSo"}
	fmt.Printf("Active Resources: %v \n", activeResources)
	for _, kind := range kindList {
		switch kind {
		case "Service":
			list := &corev1.ServiceList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeResources, rsrc.GetName()) {
					fmt.Printf("Deleting Service %v:%v \n", i, rsrc.GetName())
					err := client.Delete(ctx, &rsrc)
					if err != nil {
						fmt.Printf("Could not delete Service: %v \n", err)
					}
				}
			}
		case "Volume":
			list := &corev1.PersistentVolumeClaimList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeResources, rsrc.GetName()) {
					fmt.Printf("Deleting PVC %v: %v \n", i, rsrc.GetName())
					err := client.Delete(ctx, &rsrc)
					if err != nil {
						fmt.Printf("Could not delete PVC: %v \n", err)
					}
				}
			}
		case "Secret":
			list := &corev1.SecretList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeResources, rsrc.GetName()) {
					fmt.Printf("Deleting Secret %v: %v \n", i, rsrc.GetName())
					err := client.Delete(ctx, &rsrc)
					if err != nil {
						fmt.Printf("Could not delete Secret: %v \n", err)
					}
				}
			}
		case "Deployment":
			list := &appsv1.DeploymentList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeResources, rsrc.GetName()) {
					fmt.Printf("Deleting Deployment %v: %v \n", i, rsrc.GetName())
					err := client.Delete(ctx, &rsrc)
					if err != nil {
						fmt.Printf("Could not delete Deployment: %v \n", err)
					}
				}
			}
		case "HttpSo":
			list := &kedahttpv1alpha1.HTTPScaledObjectList{}
			err := client.List(ctx, list, &k8cl.ListOptions{Namespace: namespace})
			if err != nil {
				return err
			}
			for i, rsrc := range list.Items {
				if !slices.Contains(activeResources, rsrc.GetName()) {
					fmt.Printf("Deleting HttpSo %v: %v \n", i, rsrc.GetName())
					err := client.Delete(ctx, &rsrc)
					if err != nil {
						fmt.Printf("Could not delete HttpSo: %v \n", err)
					}
				}
			}
		}
	}
	return nil
}
