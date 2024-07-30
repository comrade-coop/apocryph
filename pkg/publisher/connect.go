// SPDX-License-Identifier: GPL-3.0

package publisher

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
	"net/url"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/comrade-coop/apocryph/pkg/abi"
	"github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/libp2p/go-libp2p/core/peer"

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

	client := pbcon.NewProvisionPodServiceClient(
		http.DefaultClient,
		url.String(),
		connect.WithInterceptors(interceptor),
	)

	return &P2pProvisionPodServiceClient{
		addr,
		client,
	}, nil
}

func SendToProvider(ctx context.Context, ipfsP2p *ipfs.P2pApi, pod *pb.Pod, deployment *pb.Deployment, client *P2pProvisionPodServiceClient, ethClient *ethclient.Client, publisherAuth *bind.TransactOpts) error {
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
				return fmt.Errorf("Failed executing provision pod request: %w", err)
			}
		} else {
			request := &pb.UpdatePodRequest{
				Pod: pod,
			}

			response, err = client.UpdatePod(ctx, connect.NewRequest(request))
			if err != nil {
				return fmt.Errorf("Failed executing update pod request: %w", err)
			}
		}

		if response.Msg.Error != "" {
			return fmt.Errorf("Error from provider: %w", errors.New(response.Msg.Error))
		}

		if response.Msg.PubAddress != "" {
			// authorize the public address to control the payment channel
			// get a payment contract instance
			payment, err := abi.NewPayment(common.Address(deployment.Payment.PaymentContractAddress), ethClient)
			if err != nil {
				return fmt.Errorf("Failed instantiating payment contract: %w", err)
			}
			tx, err := payment.Authorize(publisherAuth, common.HexToAddress(response.Msg.PubAddress), common.Address(deployment.Provider.EthereumAddress), [32]byte(deployment.Payment.PodID))
			if err != nil {
				return fmt.Errorf("Failed Authorizing Address: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Authorized Address Successfully %v\n", tx)
		}

		deployment.Deployed = response.Msg
		fmt.Fprintf(os.Stdout, "Successfully deployed! %v\n", response)
	} else {
		request := &pb.DeletePodRequest{}
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
