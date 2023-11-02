package kubernetes

import (
	"strings"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"google.golang.org/protobuf/encoding/protojson"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LabelTrustedPodsNamespace            string = "coop.comrade/trusted-pods-namespace"
	AnnotationsTrustedPodsPaymentChannel string = "coop.comrade/trusted-pods-payment-contract"
	LabelIpfsP2P                         string = "coop.comrade/trusted-pods-p2p-helper"
	AnnotationsIpfsP2P                   string = "coop.comrade/trusted-pods-p2p-helper"
)

var TrustedPodsNamespaceFilter = client.HasLabels{LabelTrustedPodsNamespace}

func NewTrustedPodsNamespace(paymentChannel *pb.PaymentChannel) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tpod-" + strings.ToLower(string(paymentChannel.PublisherAddress)),
			Labels: map[string]string{
				LabelTrustedPodsNamespace: "true",
				"pubkey":                  string(paymentChannel.PublisherAddress),
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
