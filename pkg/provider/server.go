// SPDX-License-Identifier: GPL-3.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/comrade-coop/trusted-pods/pkg/loki"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	pbcon "github.com/comrade-coop/trusted-pods/pkg/proto/protoconnect"
	"github.com/ipfs/kubo/client/rpc"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type provisionPodServer struct {
	pbcon.UnimplementedProvisionPodServiceHandler
	ipfs             *rpc.HttpApi
	k8cl             client.Client
	loki             loki.LokiConfig
	paymentValidator *ethereum.PaymentChannelValidator
	localOciRegistry string
	dryRun           bool
}

func transformError(err error) (*connect.Response[pb.ProvisionPodResponse], error) {
	return connect.NewResponse(&pb.ProvisionPodResponse{
		Error: err.Error(),
	}), nil
}

func (s *provisionPodServer) DeletePod(ctx context.Context, request *connect.Request[pb.DeletePodRequest]) (*connect.Response[pb.DeletePodResponse], error) {
	log.Println("Received request for pod deletion")

	namespace := pbcon.GetNamespace(request)
	// Create a new namespace object
	ns := &v1.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name:   namespace,
			Labels: map[string]string{},
		},
	}

	err := s.k8cl.Delete(ctx, ns)
	if err != nil {
		log.Printf("Could not delete namespace: %v\n", request)
		return nil, err
	}
	response := &pb.DeletePodResponse{Success: true}
	return connect.NewResponse(response), nil
}

func (s *provisionPodServer) UpdatePod(ctx context.Context, request *connect.Request[pb.UpdatePodRequest]) (*connect.Response[pb.ProvisionPodResponse], error) {
	log.Println("Received request for updating pod")
	secrets, err := DownloadSecrets(ctx, s.ipfs, request.Msg.Pod)
	if err != nil {
		return transformError(err)
	}
	images, err := DownloadImages(ctx, s.ipfs, s.localOciRegistry, request.Msg.Pod)
	if err != nil {
		return transformError(err)
	}

	namespace := pbcon.GetNamespace(request)
	response := &pb.ProvisionPodResponse{}
	err = tpk8s.ApplyPodRequest(ctx, s.k8cl, namespace, true, request.Msg.Pod, images, secrets, response)
	if err != nil {
		return transformError(err)
	}

	return connect.NewResponse(response), nil
}

func (s *provisionPodServer) GetPodLogs(ctx context.Context, request *connect.Request[pb.PodLogRequest], srv *connect.ServerStream[pb.PodLogResponse]) error {
	log.Println("Received Log request")

	namespace := pbcon.GetNamespace(request)
	podId := strings.Split(namespace, "-")[1]
	deploymentName := "tpod-dep-" + podId
	deployment := appsv1.Deployment{}
	key := client.ObjectKey{Namespace: namespace, Name: deploymentName}
	err := s.k8cl.Get(context.Background(), key, &deployment)
	if err != nil {
		return err
	}
	for key, value := range deployment.GetLabels() {
		if key == "containers" {
			containers := strings.Split(value, "_")
			if !slices.Contains(containers, request.Msg.ContainerName) {
				return errors.New("Container Does not exist")
			}
		}
	}

	err = loki.GetStreamedEntries(namespace, request.Msg.ContainerName, srv, s.loki.Host)
	if err != nil {
		return err
	}
	log.Println("Finished Streaming logs")
	return nil
}

func (s *provisionPodServer) ProvisionPod(ctx context.Context, request *connect.Request[pb.ProvisionPodRequest]) (*connect.Response[pb.ProvisionPodResponse], error) {
	fmt.Println("Received request for pod deployment")
	namespace := pbcon.GetNamespace(request)

	// TODO should return error (just useful for now in testing lifecycle without payment)
	if s.paymentValidator != nil {
		_, err := s.paymentValidator.Parse(request.Msg.Payment)
		if err != nil {
			return transformError(err)
		}
	}

	secrets, err := DownloadSecrets(ctx, s.ipfs, request.Msg.Pod)
	if err != nil {
		return transformError(err)
	}
	images, err := DownloadImages(ctx, s.ipfs, s.localOciRegistry, request.Msg.Pod)
	if err != nil {
		return transformError(err)
	}

	response := &pb.ProvisionPodResponse{}
	ns := tpk8s.NewTrustedPodsNamespace(namespace, request.Msg.Payment)
	err = tpk8s.RunInNamespaceOrRevert(ctx, s.k8cl, ns, s.dryRun, func(cl client.Client) error {
		return tpk8s.ApplyPodRequest(ctx, cl, ns.ObjectMeta.Name, false, request.Msg.Pod, images, secrets, response)
	})
	if err != nil {
		return transformError(err)
	}
	response.Namespace = namespace

	fmt.Println("Request processed successfully")

	return connect.NewResponse(response), nil
}

func NewTPodServerHandler(ipfsApi *rpc.HttpApi, dryRun bool, k8cl client.Client, localOciRegistry string, validator *ethereum.PaymentChannelValidator, lokiHost string) (string, http.Handler) {
	return pbcon.NewProvisionPodServiceHandler(&provisionPodServer{
		ipfs:             ipfsApi,
		k8cl:             k8cl,
		loki:             loki.LokiConfig{Host: lokiHost, Limit: "100"},
		paymentValidator: validator,
		localOciRegistry: localOciRegistry,
		dryRun:           dryRun,
	}, connect.WithInterceptors(
		pbcon.NewAuthInterceptor(k8cl),
	))
}

func GetListener(serveAddress string) (net.Listener, error) {
	listener, err := net.Listen("tcp", serveAddress)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
