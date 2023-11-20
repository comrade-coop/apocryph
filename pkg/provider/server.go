package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	"github.com/comrade-coop/trusted-pods/pkg/loki"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/gogo/status"
	"github.com/ipfs/kubo/client/rpc"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	appsv1 "k8s.io/api/apps/v1"
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

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
	}
	namespace := md.Get("namespace")[0]
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
	secrets, err := DownloadSecrets(ctx, s.ipfs, request.Pod)
	if err != nil {
		return transformError(err)
	}
	images, err := DownloadImages(ctx, s.ipfs, s.localOciRegistry, request.Pod)
	if err != nil {
		return transformError(err)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
	}
	namespace := md.Get("namespace")[0]
	response := &pb.ProvisionPodResponse{}
	err = tpk8s.ApplyPodRequest(ctx, s.k8cl, namespace, true, request.Pod, images, secrets, response)

	return response, err
}

func (s *provisionPodServer) GetPodLogs(request *pb.PodLogRequest, srv pb.ProvisionPodService_GetPodLogsServer) error {
	log.Println("Received Log request")

	md, ok := metadata.FromIncomingContext(srv.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata not found")
	}

	namespace := md.Get("namespace")[0]
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
			if !slices.Contains(containers, request.ContainerName) {
				return errors.New("Container Does not exist")
			}
		}
	}

	err = loki.GetStreamedEntries(namespace, request.ContainerName, srv, s.loki.Host)
	if err != nil {
		return err
	}
	log.Println("Finished Streaming logs")
	return nil
}

func (s *provisionPodServer) ProvisionPod(ctx context.Context, request *pb.ProvisionPodRequest) (*pb.ProvisionPodResponse, error) {
	fmt.Println("Received request for pod deployment")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
	}
	namespace := md.Get("namespace")[0]

	// TODO should return error (just usefull for now in testing lifecycle without payment)
	if s.paymentValidator != nil {
		_, err := s.paymentValidator.Parse(request.Payment)
		if err != nil {
			return transformError(err)
		}
	}

	secrets, err := DownloadSecrets(ctx, s.ipfs, request.Pod)
	if err != nil {
		return transformError(err)
	}
	images, err := DownloadImages(ctx, s.ipfs, s.localOciRegistry, request.Pod)
	if err != nil {
		return transformError(err)
	}

	response := &pb.ProvisionPodResponse{}
	ns := tpk8s.NewTrustedPodsNamespace(namespace, request.Payment)
	err = tpk8s.RunInNamespaceOrRevert(ctx, s.k8cl, ns, s.dryRun, func(cl client.Client) error {
		return tpk8s.ApplyPodRequest(ctx, cl, ns.ObjectMeta.Name, false, request.Pod, images, secrets, response)
	})
	if err != nil {
		return transformError(err)
	}
	response.Namespace = namespace

	fmt.Println("Request processed successfully")

	return response, nil
}

func NewTPodServer(ipfsApi *rpc.HttpApi, dryRun bool, k8cl client.Client, localOciRegistry string, validator *ethereum.PaymentChannelValidator, lokiHost string) (*grpc.Server, error) {
	server := grpc.NewServer(
		grpc.ChainStreamInterceptor(pb.NoCrashStreamServerInterceptor),
		grpc.ChainStreamInterceptor(pb.NoCrashStreamServerInterceptor, pb.AuthStreamServerInterceptor(k8cl)),
		grpc.ChainUnaryInterceptor(pb.NoCrashUnaryServerInterceptor, pb.AuthUnaryServerInterceptor(k8cl)),
	)

	pb.RegisterProvisionPodServiceServer(server, &provisionPodServer{
		ipfs:             ipfsApi,
		k8cl:             k8cl,
		loki:             loki.LokiConfig{Host: lokiHost, Limit: "100"},
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
