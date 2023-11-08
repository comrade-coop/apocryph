package provider

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/comrade-coop/trusted-pods/pkg/loki"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/kubo/client/rpc"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type provisionPodServer struct {
	pb.UnimplementedProvisionPodServiceServer
	ipfs             *rpc.HttpApi
	k8cl             client.Client
	loki             loki.LokiConfig
	paymentValidator *ethereum.PaymentChannelValidator
	localOciRegistry string
	dryRun           bool
}

func transformError(err error) (*pb.ProvisionPodResponse, error) {
	return &pb.ProvisionPodResponse{
		Error: err.Error(),
	}, nil
}

func (s *provisionPodServer) DeletePod(ctx context.Context, request *pb.DeletePodRequest) (*pb.DeletePodResponse, error) {
	log.Println("Received request for pod deletion")

	namespace := "tpod-" + strings.ToLower(common.BytesToAddress(request.Credentials.PublisherAddress).String())
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
	return response, nil
}

func (s *provisionPodServer) UpdatePod(ctx context.Context, request *pb.UpdatePodRequest) (*pb.ProvisionPodResponse, error) {
	log.Println("Received request for updating pod")
	err := tpipfs.TransformSecrets(request.Pod, tpipfs.DownloadSecrets(ctx, s.ipfs), tpipfs.DecryptSecrets(request.Keys))
	if err != nil {
		return nil, err
	}

	if s.localOciRegistry != "" {
		err = tpipfs.ReuploadImagesFromIpdr(request.Pod, ctx, s.ipfs, s.localOciRegistry, nil, request.Keys)
		if err != nil {
			return nil, err
		}
	}

	namespace := "tpod-" + strings.ToLower(common.BytesToAddress(request.Credentials.PublisherAddress).String())

	response := &pb.ProvisionPodResponse{}
	err = tpk8s.ApplyPodRequest(ctx, s.k8cl, request.Pod, true, namespace, response)

	return response, err
}

func (s *provisionPodServer) GetPodLogs(request *pb.PodLogRequest, srv pb.ProvisionPodService_GetPodLogsServer) error {
	log.Println("Received Log request")
	// verify if container exists or not
	p := &v1.Namespace{}

	namespace := "tpod-" + strings.ToLower(common.BytesToAddress(request.Credentials.PublisherAddress).String())
	err := s.k8cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: namespace}, p)
	if err != nil {
		return err
	}

	err = loki.GetStreamedEntries(namespace, request.ContainerName, srv)
	if err != nil {
		return err
	}
	log.Println("Finished Streaming logs")
	return nil
}

func (s *provisionPodServer) ProvisionPod(ctx context.Context, request *pb.ProvisionPodRequest) (*pb.ProvisionPodResponse, error) {
	fmt.Println("Received request for pod deployment")
	pod := request.Pod

	var err error

	// TODO should return error (just usefull for now in testing lifecycle without payment)
	if s.paymentValidator != nil {
		_, err = s.paymentValidator.Parse(request.Payment)
		if err != nil {
			return transformError(err)
		}
	}

	err = tpipfs.TransformSecrets(pod, tpipfs.DownloadSecrets(ctx, s.ipfs), tpipfs.DecryptSecrets(request.Keys))
	if err != nil {
		return transformError(err)
	}

	if s.localOciRegistry != "" {
		err = tpipfs.ReuploadImagesFromIpdr(pod, ctx, s.ipfs, s.localOciRegistry, nil, request.Keys)
		if err != nil {
			return transformError(err)
		}
	}

	response := &pb.ProvisionPodResponse{}
	namespace := tpk8s.NewTrustedPodsNamespace(request.Payment)
	err = tpk8s.RunInNamespaceOrRevert(ctx, s.k8cl, namespace, s.dryRun, func(cl client.Client) error {
		return tpk8s.ApplyPodRequest(ctx, cl, pod, false, namespace.ObjectMeta.Name, response)
	})
	if err != nil {
		return transformError(err)
	}
	response.Namespace = namespace.GetName()

	fmt.Println("Request processed successfully")

	return response, nil
}

func NewTPodServer(ipfsApi *rpc.HttpApi, dryRun bool, k8cl client.Client, localOciRegistry string, validator *ethereum.PaymentChannelValidator) (*grpc.Server, error) {
	server := grpc.NewServer(grpc.UnaryInterceptor(pb.AuthUnaryServerInterceptor()))

	pb.RegisterProvisionPodServiceServer(server, &provisionPodServer{
		ipfs:             ipfsApi,
		k8cl:             k8cl,
		loki:             loki.LokiConfig{Url: "http://loki.loki.svc.cluster.local:3100/loki", Limit: "100"},
		paymentValidator: validator,
		localOciRegistry: localOciRegistry,
		dryRun:           dryRun,
	})
	return server, nil
}

func GetListener(serveAddress string) (net.Listener, error) {
	listener, err := net.Listen("tcp", serveAddress)
	if err != nil {
		return nil, err
	}
	return listener, nil
}
