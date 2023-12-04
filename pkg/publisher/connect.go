package publisher

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"connectrpc.com/connect"
	"github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/provider"
	"github.com/ethereum/go-ethereum/common"
	pbcon "github.com/comrade-coop/trusted-pods/pkg/proto/protoconnect"
	"github.com/libp2p/go-libp2p/core/peer"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
)

type P2pProvisionPodServiceClient struct {
	*tpipfs.IpfsAddr
	pbcon.ProvisionPodServiceClient
}

func ConnectToProvider(ipfsP2p *ipfs.P2pApi, deployment *pb.Deployment, interceptor *pbcon.AuthInterceptorClient, availableProviders []*provider.HostInfo) (P2pProvisionPodServiceClient, error) {
	deployment.Provider = &pb.ProviderConfig{}
	for _, provider := range availableProviders {
		deployment.Provider.Libp2PAddress = provider.MultiAddresses[0].Value
		deployment.Provider.EthereumAddress = common.HexToAddress(provider.Id).Bytes()
		providerPeerId, err := peer.Decode(deployment.GetProvider().GetLibp2PAddress())
		if err != nil {
			return P2pProvisionPodServiceClient{}, fmt.Errorf("Failed to parse provider address: %w", err)
		}
		// TODO add ping protocol to test connection before deployment
		addr, err := ipfsP2p.ConnectTo(pb.ProvisionPod, providerPeerId)
		if err != nil {
			//return P2pProvisionPodServiceClient{}, fmt.Errorf("Failed to dial provider: %w", err)
			continue
		}

		url := &url.URL{Scheme: "http", Host: addr.String()}

		client := pbcon.NewProvisionPodServiceClient(
			http.DefaultClient,
			url.String(),
			connect.WithInterceptors(interceptor),
		)

		return P2pProvisionPodServiceClient{
			addr,
			client,
		}, nil
	}
	return P2pProvisionPodServiceClient{}, fmt.Errorf("Failed to dial provider(s)")
}

func SendToProvider(ctx context.Context,
	ipfsP2p *ipfs.P2pApi, pod *pb.Pod, deployment *pb.Deployment, interceptor *pbcon.AuthInterceptorClient,
	availableProviders []*provider.HostInfo, fundPaymentChannelFunc func() error) error {
	// tpipfs.NewP2pApi(ipfs, ipfsMultiaddr)
	pod = LinkUploadsFromDeployment(pod, deployment)

	client, err := ConnectToProvider(ipfsP2p, deployment, interceptor, availableProviders)
	if err != nil {
		return err
	}
	defer client.Close()

	if fundPaymentChannelFunc != nil {
		err = fundPaymentChannelFunc()
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stderr, "Sending request to provider over IPFS p2p...\n")

	if pod != nil {
		var response *connect.Response[pb.ProvisionPodResponse]
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
			response, err = client.ProvisionPod(ctx, connect.NewRequest(request))
			if err != nil {
				return fmt.Errorf("Failed executing provision pod request: %w", err)
			}
		} else {
			request := &pb.UpdatePodRequest{
				Pod:              pod,
				PublisherAddress: deployment.Payment.PublisherAddress,
			}
			response, err = client.UpdatePod(ctx, connect.NewRequest(request))
			if err != nil {
				return fmt.Errorf("Failed executing update pod request: %w", err)
			}
		}

		if response.Msg.Error != "" {
			return fmt.Errorf("Error from provider: %w", errors.New(response.Msg.Error))
		}

		deployment.Deployed = response.Msg
		fmt.Fprintf(os.Stderr, "Successfully deployed! %v\n", response.Msg)
	} else {
		request := &pb.DeletePodRequest{
			PublisherAddress: deployment.Payment.PublisherAddress,
		}
		response, err := client.DeletePod(ctx, connect.NewRequest(request))
		if err != nil {
			return fmt.Errorf("Failed executing update pod request: %w", err)
		}
		if response.Msg.Error != "" {
			return fmt.Errorf("Error from provider: %w", errors.New(response.Msg.Error))
		}
		deployment.Deployed = nil
		fmt.Fprintf(os.Stderr, "Successfully undeployed!\n")
	}

	return nil
}
