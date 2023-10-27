package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type provisionPodServer struct {
	pb.UnimplementedProvisionPodServiceServer
	ipfs             *rpc.HttpApi
	k8cl             client.Client
	paymentValidator *ethereum.PaymentChannelValidator
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
func (s *provisionPodServer) UpdatePod(ctx context.Context, request *pb.UpdatePodRequest) (*pb.UpdatePodResponse, error) {
	log.Println("Received request for updating pod")

	serviceList := &v1.ServiceList{}
	deploymentList := &apps.DeploymentList{}

	listOptions := []client.ListOption{
		client.InNamespace(request.Namespace),
	}
	if len(request.Pod.Containers) != 0 {
		for cIdx, container := range request.Pod.Containers {
			log.Printf("Processing container: %v\n", container)
			if len(container.Ports) != 0 {
				var newService *v1.Service
				var err error
				for i, port := range container.Ports {
					log.Printf("Processing port: %v\n", port)
					portName := fmt.Sprintf("p%d-%d", cIdx, port.ContainerPort)
					httpSo := tpk8s.CreateHttpSo()
					labels := map[string]string{"tpod": "1"}
					newService, _, err = tpk8s.GetService(port, portName, httpSo, labels)
					if err != nil {
						return nil, err
					}
					err := s.k8cl.List(ctx, serviceList, listOptions...)
					if err != nil {
						log.Printf("retreiving services for namespace failed: %v\n", err)
						return nil, err
					}
					service := serviceList.Items[i]
					patch, err := json.Marshal(newService)
					if err != nil {
						log.Printf("Marshall faillure: %v\n", err)
						return nil, err
					}
					err = s.k8cl.Patch(ctx, &service, client.RawPatch(types.MergePatchType, patch))
					if err != nil {
						log.Printf("Failed patching service: %v\n", err)
						return nil, err
					}
				}
			}
		}
	}
	err := s.k8cl.List(ctx, deploymentList, listOptions...)
	if err != nil {
		log.Printf("retreiving deployments for namespace failed: %v\n", err)
		return nil, err
	}

	for _, obj := range deploymentList.Items {
		log.Printf("updating object %v", obj.GetName())
	}

	response := &pb.UpdatePodResponse{}

	return response, nil
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

	err = tpipfs.TransformSecrets(pod, tpipfs.DownloadSecrets(ctx, s.ipfs), tpipfs.DecryptSecrets(request.Keys))
	if err != nil {
		return transformError(err)
	}

	if s.paymentValidator != nil {
		_, err = s.paymentValidator.Parse(request.Payment)
		if err != nil {
			return transformError(err)
		}

	}

	response := &pb.ProvisionPodResponse{}
	namespace := tpk8s.NewTrustedPodsNamespace(request.Payment)
	err = tpk8s.RunInNamespaceOrRevert(ctx, s.k8cl, namespace, s.dryRun, func(cl client.Client) error {
		return tpk8s.ApplyPodRequest(ctx, cl, pod, response)
	})
	if err != nil {
		return transformError(err)
	}
	response.Namespace = namespace.GetName()

	fmt.Printf("Request processed successfully, %v %v\n", response, namespace)

	return response, nil
}

func InitTPodServer(ipfsApi *rpc.HttpApi, kubeConfig string, dryRun bool, val ...*ethereum.PaymentChannelValidator) (*grpc.Server, error) {
	server := grpc.NewServer()

	k8cl, err := tpk8s.GetClient(kubeConfig, dryRun)
	if err != nil {
		return nil, err
	}
	var validator *ethereum.PaymentChannelValidator
	if len(val) > 0 {
		validator = val[0]
	} else {
		validator = nil
	}
	pb.RegisterProvisionPodServiceServer(server, &provisionPodServer{
		ipfs:             ipfsApi,
		k8cl:             k8cl,
		paymentValidator: validator,
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
