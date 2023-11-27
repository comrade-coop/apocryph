package publisher

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/provider"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
)

func ConnectToProvider(ipfsP2p *ipfs.P2pApi, deployment *pb.Deployment, interceptor *pb.AuthInterceptorClient, availableProviders []*provider.HostInfo) (*tpipfs.IpfsClientConn, error) {
	for _, provider := range availableProviders {
		deployment.Provider.Libp2PAddress = provider.Id
		providerPeerId, err := peer.Decode(deployment.GetProvider().GetLibp2PAddress())
		if err != nil {
			return nil, fmt.Errorf("Failed to parse provider address: %w", err)
		}
		conn, err := ipfsP2p.ConnectToGrpc(
			pb.ProvisionPod,
			providerPeerId,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(interceptor.StreamClientInterceptor()),
		)
		if err != nil {
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("Failed to dial provider(s)")
}

func SendToProvider(ctx context.Context, ipfsP2p *ipfs.P2pApi, pod *pb.Pod, deployment *pb.Deployment, interceptor *pb.AuthInterceptorClient, availableProviders []*provider.HostInfo) error {
	// tpipfs.NewP2pApi(ipfs, ipfsMultiaddr)
	pod = LinkUploadsFromDeployment(pod, deployment)

	conn, err := ConnectToProvider(ipfsP2p, deployment, interceptor, availableProviders)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewProvisionPodServiceClient(conn)

	fmt.Fprintf(os.Stderr, "Sending request to provider over IPFS p2p...\n")

	if pod != nil {
		var response *pb.ProvisionPodResponse
		if deployment.Deployed == nil || deployment.Deployed.Error != "" {
			request := &pb.ProvisionPodRequest{
				Pod: pod,
				Payment: &pb.PaymentChannel{
					ChainID:          deployment.Payment.ChainID,
					ProviderAddress:  deployment.Provider.EthereumAddress,
					ContractAddress:  deployment.Payment.PaymentContractAddress,
					PublisherAddress: deployment.Payment.PublisherAddress,
					PodID:            deployment.Payment.PodID,
				},
			}
			response, err = client.ProvisionPod(ctx, request)
			if err != nil {
				return fmt.Errorf("Failed executing provision pod request: %w", err)
			}
		} else {
			request := &pb.UpdatePodRequest{
				Pod:              pod,
				PublisherAddress: deployment.Payment.PublisherAddress,
			}
			response, err = client.UpdatePod(ctx, request)
			if err != nil {
				return fmt.Errorf("Failed executing update pod request: %w", err)
			}
		}

		if response.Error != "" {
			return fmt.Errorf("Error from provider: %w", errors.New(response.Error))
		}

		deployment.Deployed = response
		fmt.Fprintf(os.Stderr, "Successfully deployed! %v\n", response)
	} else {
		request := &pb.DeletePodRequest{
			PublisherAddress: deployment.Payment.PublisherAddress,
		}
		response, err := client.DeletePod(ctx, request)
		if err != nil {
			return fmt.Errorf("Failed executing update pod request: %w", err)
		}
		if response.Error != "" {
			return fmt.Errorf("Error from provider: %w", errors.New(response.Error))
		}
		deployment.Deployed = nil
		fmt.Fprintf(os.Stderr, "Successfully undeployed!\n")
	}

	return nil
}
