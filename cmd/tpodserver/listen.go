package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"time"

	ipfs_utils "github.com/comrade-coop/trusted-pods/pkg/ipfs-utils"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var ipfsApi string
var serveAddress string

type provisionPodServer struct {
	pb.UnimplementedProvisionPodServiceServer
	ipfs *rpc.HttpApi
	k8cl client.Client
}

func transformError(err error) (*pb.ProvisionPodResponse, error) {
	return &pb.ProvisionPodResponse{
		Error: err.Error(),
	}, nil
}

func VerifyPaymentChannel(client common.Address, podID *big.Int, token common.Address) error {

	channel, err := Instance.Channels(&bind.CallOpts{Pending: false}, client, common.HexToAddress(ProviderAddress), podID, token)
	if err != nil {
		return err
	}
	// verify if deadline reached
	if channel.Deadline.Int64() < time.Now().Unix() {
		return errors.New("channel expired")
	}
	if channel.MinAdvanceDuration.Int64() < MinExecutionPeriod {
		return errors.New("Min deadline prolonging Period is less than MinExecutionPeriod")
	}
	if channel.Price.Int64() < PricePerExecution {
		return errors.New("Channel Price is less than PricePerExecution")
	}
	return nil
}

func (s *provisionPodServer) ProvisionPod(ctx context.Context, request *pb.ProvisionPodRequest) (*pb.ProvisionPodResponse, error) {
	fmt.Printf("Received request for pod deployment, %v", request)
	cid, err := cid.Cast(request.PodManifestCid)
	if err != nil {
		return transformError(err)
	}
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

	err = VerifyPaymentChannel(common.HexToAddress(request.ClientAddress), new(big.Int).SetBytes(request.PodID), common.HexToAddress(request.TokenAddress))
	if err != nil {
		return transformError(err)
	}

	response := &pb.ProvisionPodResponse{}
	namespace := tpk8s.NewTrustedPodsNamespace(paymentContract)
	err = tpk8s.RunInNamespaceOrRevert(ctx, s.k8cl, namespace, dryRun, func(cl client.Client) error {
		return tpk8s.ApplyPodRequest(ctx, cl, s.ipfs, request.Keys, pod, response)
	})
	if err != nil {
		return transformError(err)
	}

	fmt.Printf("Request processed successfully, %v %v", response, namespace)

	return response, nil
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start a service listening for incomming execution requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := setUp()
		if err != nil {
			return err
		}

		ipfs, ipfsMultiaddr, err := ipfs_utils.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		k8cl, err := tpk8s.GetClient(kubeConfig, dryRun)
		if err != nil {
			return err
		}

		listener, err := GetListener(ipfs, ipfsMultiaddr, serveAddress)
		if err != nil {
			return err
		}

		server := grpc.NewServer()
		pb.RegisterProvisionPodServiceServer(server, &provisionPodServer{
			ipfs: ipfs,
			k8cl: k8cl,
		})

		go server.Serve(listener)

		defer server.Stop()

		<-cmd.Context().Done()

		server.GracefulStop()

		return nil
	},
}

func GetListener(ipfs *rpc.HttpApi, ipfsMultiaddr multiaddr.Multiaddr, serveAddress string) (net.Listener, error) {
	if serveAddress == "-" {
		return ipfs_utils.NewP2pApi(ipfs, ipfsMultiaddr).Listen(pb.ProvisionPod)
	} else {
		return net.Listen("tcp", serveAddress)
	}
}

func init() {
	listenCmd.Flags().StringVar(&serveAddress, "address", "-", "port to serve on (- to automatically pick a port and register a listener for it in ipfs)")

	listenCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	listenCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "-", "absolute path to the kubeconfig file (- to the first of in-cluster config and ~/.kube/config)")
	listenCmd.Flags().StringVar(&ipfsApi, "ipfs", "-", "multiaddr where the ipfs/kubo api can be accessed (- to use the daemon running in IPFS_PATH)")
}
