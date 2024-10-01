// SPDX-License-Identifier: GPL-3.0

package publisher

import (
	"context"
	"errors"
	"fmt"
	"net/rpc"
	"net/url"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
	"github.com/libp2p/go-libp2p/core/peer"
	retryablehttp  "github.com/hashicorp/go-retryablehttp"

	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
)

type P2pProvisionPodServiceClient struct {
	*tpipfs.IpfsAddr
	pbcon.ProvisionPodServiceClient
}

func ConnectToProvider(ipfsP2p *ipfs.P2pApi, deployment *pb.Deployment, interceptor *pbcon.AuthInterceptorClient) (*P2pProvisionPodServiceClient, error) {
	providerPeerId, err := peer.Decode(strings.TrimPrefix(deployment.GetProvider().GetLibp2PAddress(), "/p2p/"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse provider address: %w", err)
	}
	addr, err := ipfsP2p.ConnectTo(pb.ProvisionPod, providerPeerId)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial provider: %w", err)
	}

	url := &url.URL{Scheme: "http", Host: addr.String()}

	// ping the provider
	pingClient, err := rpc.Dial("tcp", url.Host)
	if err != nil {
		return nil, fmt.Errorf("Error pinging %s: %s\n", url, err)
	}
	pingClient.Close()

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 8
	
	client := pbcon.NewProvisionPodServiceClient(
		retryClient.StandardClient(),
		url.String(),
		connect.WithInterceptors(interceptor),
	)

	return &P2pProvisionPodServiceClient{
		addr,
		client,
	}, nil
}

func SendToProvider(ctx context.Context, ipfsP2p *ipfs.P2pApi, pod *pb.Pod, deployment *pb.Deployment, client *P2pProvisionPodServiceClient) (interface{}, error) {
	// tpipfs.NewP2pApi(ipfs, ipfsMultiaddr)
	pod = LinkUploadsFromDeployment(pod, deployment)
	defer client.Close()
	fmt.Fprintf(os.Stderr, "Sending request to provider over IPFS p2p...\n")

	if pod != nil {
		var err error
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
			fmt.Println("Processing Request ...")
			response, err = client.ProvisionPod(ctx, connect.NewRequest(request))
			if err != nil {
				return nil, fmt.Errorf("Failed executing provision pod request: %w", err)
			}
		} else {
			request := &pb.UpdatePodRequest{
				Pod: pod,
				Payment: &pb.PaymentChannel{
					ChainID:          deployment.Payment.ChainID,
					ProviderAddress:  deployment.Provider.EthereumAddress,
					ContractAddress:  deployment.Payment.PaymentContractAddress,
					PublisherAddress: deployment.Payment.PublisherAddress,
					PodID:            deployment.Payment.PodID,
				},
			}

			response, err = client.UpdatePod(ctx, connect.NewRequest(request))
			if err != nil {
				return nil, fmt.Errorf("Failed executing update pod request: %w", err)
			}
			return response.Msg, nil
		}

		if response.Msg.Error != "" {
			return nil, fmt.Errorf("Error from provider: %w", errors.New(response.Msg.Error))
		}
		deployment.Deployed = response.Msg
		fmt.Fprintf(os.Stdout, "Successfully deployed! %v\n", response)

		return response.Msg, nil

	} else {
		request := &pb.DeletePodRequest{}
		response, err := client.DeletePod(ctx, connect.NewRequest(request))
		if err != nil {
			return nil, fmt.Errorf("Failed executing delete pod request: %w", err)
		}
		if response.Msg.Error != "" {
			return nil, fmt.Errorf("Error from provider: %w", errors.New(response.Msg.Error))
		}
		deployment.Deployed = nil
		fmt.Fprintf(os.Stderr, "Successfully undeployed!\n")
		return response.Msg, nil
	}
}
