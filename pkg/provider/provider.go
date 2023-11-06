package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

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
	log.Println("Received request for pod deletion")

	namespace := "tpod-" + strings.ToLower(string(request.Credentials.PublisherAddress))
	// Create a new namespace object
	ns := &v1.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"pubkey": string(request.Credentials.PublisherAddress),
			},
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

	namespace := "tpod-" + strings.ToLower(string(request.Credentials.PublisherAddress))

	response := &pb.ProvisionPodResponse{}
	err = tpk8s.ApplyPodRequest(ctx, s.k8cl, request.Pod, true, namespace, response)

	return response, err
}

func (s *provisionPodServer) GetPodLogs(request *pb.PodLogRequest, srv pb.ProvisionPodService_GetPodLogsServer) error {
	log.Println("Received Log request")
	// verify if container exists or not
	namespace := "tpod-" + strings.ToLower(string(request.Credentials.PublisherAddress))
	err := s.k8cl.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: request.ContainerName}, nil)
	if err != nil {
		return errors.New("Container Does not exists")
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
	cid, err := cid.Cast(request.PodManifestCid)
	if err != nil {
		return transformError(err)
	}
	fmt.Println("Retreiving pod from Ipfs")
	// TODO verify CID size before downloading from ipfs
	node, err := s.ipfs.Unixfs().Get(ctx, path.IpfsPath(cid))
	if err != nil {
		return transformError(err)
	}
	file, ok := node.(files.File)
	if !ok {
		return transformError(errors.New("Supplied CID not a file"))
	}
	fmt.Println("Reading pod manifest")
	manifestBytes, err := io.ReadAll(file)
	if err != nil {
		return transformError(err)
	}
	pod := &pb.Pod{}
	err = proto.Unmarshal(manifestBytes, pod)
	if err != nil {
		return transformError(err)
	}

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
