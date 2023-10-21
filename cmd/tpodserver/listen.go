package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	tpk8s "github.com/comrade-coop/trusted-pods/pkg/kubernetes"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
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
	paymentValidator *ethereum.PaymentChannelValidator
}

func transformError(err error) (*pb.ProvisionPodResponse, error) {
	return &pb.ProvisionPodResponse{
		Error: err.Error(),
	}, nil
}

func (s *provisionPodServer) ProvisionPod(ctx context.Context, request *pb.ProvisionPodRequest) (*pb.ProvisionPodResponse, error) {
	fmt.Printf("Received request for pod deployment, %v\n", request)
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

	_, err = s.paymentValidator.Parse(request.Payment)
	if err != nil {
		return transformError(err)
	}

	response := &pb.ProvisionPodResponse{}
	namespace := tpk8s.NewTrustedPodsNamespace(request.Payment)
	err = tpk8s.RunInNamespaceOrRevert(ctx, s.k8cl, namespace, dryRun, func(cl client.Client) error {
		return tpk8s.ApplyPodRequest(ctx, cl, s.ipfs, request.Keys, pod, response)
	})
	if err != nil {
		return transformError(err)
	}

	fmt.Printf("Request processed successfully, %v %v\n", response, namespace)

	return response, nil
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start a service listening for incomming execution requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		k8cl, err := tpk8s.GetClient(kubeConfig, dryRun)
		if err != nil {
			return err
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		providerAuth, err := ethereum.GetAccount(providerKey, ethClient)
		if err != nil {
			return err
		}

		pricingTable, err := openPricingTable()
		if err != nil {
			return err
		}

		validator, err := ethereum.NewPaymentChannelValidator(ethClient, allowedContractAddresses, providerAuth, pricingTable.TokenAddress)

		var listener net.Listener
		if serveAddress == "" {
			listener, err = tpipfs.NewP2pApi(ipfs, ipfsMultiaddr).Listen(pb.ProvisionPod)
		} else {
			listener, err = net.Listen("tcp", serveAddress)
		}
		if err != nil {
			return err
		}

		server := grpc.NewServer()
		pb.RegisterProvisionPodServiceServer(server, &provisionPodServer{
			ipfs: ipfs,
			k8cl: k8cl,
			paymentValidator: validator,
		})

		go server.Serve(listener)

		defer server.Stop()

		<-cmd.Context().Done()

		server.GracefulStop()

		return nil
	},
}

func init() {
	listenCmd.Flags().StringVar(&serveAddress, "address", "", "port to serve on (leave blank to automatically pick a port and register a listener for it in ipfs)")


	listenCmd.Flags().BoolVarP(&dryRun, "dry-run", "z", false, "Dry run mode; modify nothing.")
	listenCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "", "absolute path to the kubeconfig file (leave blank for the first of in-cluster config and ~/.kube/config)")
	listenCmd.Flags().StringVar(&ipfsApi, "ipfs", "", "multiaddr where the ipfs/kubo api can be accessed (leave blank to use the daemon running in IPFS_PATH)")
	listenCmd.Flags().StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "client public address")
	listenCmd.Flags().StringVar(&providerKey, "ethereum-key", "", "provider account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")
}
