package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/comrade-coop/trusted-pods/pkg/loki"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
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
	log.Printf("\n Received request for pod deletion: %v\n", request)
	// TODO Authentication
	// Create a new namespace object
	ns := &v1.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: request.Namespace,
		},
	}
	err := s.k8cl.Delete(ctx, ns)
	if err != nil {
		log.Printf("Could not delete namespace: %v\n", request)
		return nil, err
	}
	response := &pb.DeletePodResponse{Namespace: ns.GetName()}
	return response, nil
}

func (s *provisionPodServer) UpdatePod(ctx context.Context, request *pb.UpdatePodRequest) (*pb.ProvisionPodResponse, error) {
	log.Println("Received request for updating pod")

	response := &pb.ProvisionPodResponse{}
	tpk8s.ApplyPodRequest(ctx, s.k8cl, request.Pod, true, request.Namespace, response)

	return response, nil
}

func (s *provisionPodServer) GetPodLogs(ctx context.Context, request *pb.PodLogsRequest) (*pb.PodLogResponse, error) {
	log.Println("Received Log request")
	response := pb.PodLogResponse{}
	entries, err := loki.GetLogs(request.ContainerName, s.loki.Limit, s.loki.Url)
	if err != nil {
		return nil, err
	}
	response.LogEntries = entries
	return &response, nil
}

func (s *provisionPodServer) ProvisionPod(ctx context.Context, request *pb.ProvisionPodRequest) (*pb.ProvisionPodResponse, error) {
	fmt.Printf("Received request for pod deployment, %v\n", request)
	cid, err := cid.Cast(request.PodManifestCid)
	if err != nil {
		return transformError(err)
	}
	// TODO verify CID size before downloading from ipfs
	node, err := s.ipfs.Unixfs().Get(ctx, path.IpfsPath(cid))
	if err != nil {
		return transformError(err)
	}
	file, ok := node.(files.File)
	if !ok {
		return transformError(errors.New("Supplied CID not a file"))
	}
	manifestBytes, err := io.ReadAll(file)
	if err != nil {
		return transformError(err)
	}
	pod := &pb.Pod{}
	err = proto.Unmarshal(manifestBytes, pod)
	if err != nil {
		return transformError(err)
	}

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

	fmt.Printf("Request processed successfully, %v %v\n", response, namespace)

	return response, nil
}

func NewTPodServer(ipfsApi *rpc.HttpApi, dryRun bool, k8cl client.Client, localOciRegistry string, validator *ethereum.PaymentChannelValidator) (*grpc.Server, error) {
	server := grpc.NewServer()

	pb.RegisterProvisionPodServiceServer(server, &provisionPodServer{
		ipfs:             ipfsApi,
		k8cl:             k8cl,
		loki:             loki.LokiConfig{Url: "http://loki.loki.svc.cluster.local:3100/loki/api/v1/query_range", Limit: "100"},
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
