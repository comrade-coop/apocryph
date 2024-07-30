// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"math/big"
	"path/filepath"

	"github.com/comrade-coop/apocryph/pkg/abi"
	"github.com/comrade-coop/apocryph/pkg/ethereum"
	"github.com/comrade-coop/apocryph/pkg/ipcr"
	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	pbcon "github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
	"github.com/comrade-coop/apocryph/pkg/publisher"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/spf13/cobra"
)

// if no provider is selected, Fetches providers based on registry args
func fetchAndFilterProviders(ipfs *rpc.HttpApi, ethClient *ethclient.Client) (publisher.ProviderHostInfoList, error) {
	registryContract := common.HexToAddress(registryContractAddress)
	var tokenContract common.Address
	if tokenContractAddress != "" {
		tokenContract = common.HexToAddress(tokenContractAddress)
	} else if paymentContractAddress != "" {
		paymentContract := common.HexToAddress(paymentContractAddress)
		payment, err := abi.NewPayment(paymentContract, ethClient)
		if err != nil {
			return nil, fmt.Errorf("Failed instantiating payment contract: %w", err)
		}
		tokenContract, err = payment.Token(&bind.CallOpts{})
		if err != nil {
			return nil, fmt.Errorf("Failed getting payment contract token: %w", err)
		}
	} else {
		return nil, fmt.Errorf("Either a token contract or a payment contract must be set")
	}

	tables, err := publisher.GetPricingTables(ethClient, registryContract, tokenContract)
	if err != nil {
		return nil, err
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("Marketplace is empty! are you sure you are connected to the right contract?")
	}

	filter, err := getRegistryTableFilter()
	if err != nil {
		return nil, err
	}

	filteredTables := publisher.FilterPricingTables(tables, filter)
	if len(filteredTables) == 0 {
		return nil, fmt.Errorf("no table found by filter")
	}

	err = PrintTableInfo(filteredTables, ethClient, registryContract)
	if err != nil {
		return nil, err
	}

	availableProviders, err := publisher.GetProviderHostInfos(ipfs, ethClient, registryContract, filteredTables)
	if err != nil {
		return nil, err
	}

	availableProviders, err = publisher.FilterProviderHostInfos(region, providerPeer, availableProviders)
	if err != nil {
		return nil, err
	}

	PrintProvidersInfo(availableProviders)

	return availableProviders, nil
}

// FIXME: Massive code duplication with registry get, upload, fund, and sync commands.
var deployPodCmd = &cobra.Command{
	Use:     fmt.Sprintf("deploy [%s] [deployment.yaml]", publisher.DefaultPodFile),
	Aliases: []string{"update"},
	Short:   "Deploy a pod from a local manifest",
	Args:    cobra.MaximumNArgs(2),
	GroupID: "main",
	RunE: func(cmd *cobra.Command, args []string) error {
		podFile, deploymentFile, pod, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		configureDeployment(deployment)

		if authorize {
			pod.Authorized = true
		}

		fundsInt, _ := (&big.Int{}).SetString(funds, 10)
		if fundsInt == nil {
			return fmt.Errorf("Invalid number passed for funds: %s", funds)
		}

		unlockTimeInt := big.NewInt(unlockTime)

		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connecting to IPFS: %w", err)
		}

		ipfsp2p := tpipfs.NewP2pApi(ipfs, ipfsMultiaddr)

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}
		if publisherKey == "" {
			publisherKey = common.BytesToAddress(deployment.Payment.PublisherAddress).String()
		}

		publisherAuth, sign, err := ethereum.GetAccountAndSigner(publisherKey, ethClient)
		if err != nil {
			return fmt.Errorf("Could not get ethereum account: %w", err)
		}

		// FIXME: move to configureDeployment?
		chainId, err := ethClient.ChainID(cmd.Context())
		if err != nil {
			return err
		}
		deployment.Payment.ChainID = chainId.Bytes()
		deployment.Payment.PublisherAddress = publisherAuth.From.Bytes()

		interceptor := pbcon.NewAuthInterceptorClient(deployment, expirationOffset, sign)

		var provisionPodclient *publisher.P2pProvisionPodServiceClient
		if len(deployment.GetProvider().GetEthereumAddress()) == 0 || deployment.GetProvider().GetLibp2PAddress() == "" {
			availableProviders, err := fetchAndFilterProviders(ipfs, ethClient)
			if err != nil {
				return fmt.Errorf("Failed finding a provider: %w", err)
			}
			if len(availableProviders) == 0 {
				return fmt.Errorf("Failed finding a provider: no available providers found matching filter")
			}
			provisionPodclient, err = publisher.SetFirstConnectingProvider(ipfsp2p, availableProviders, deployment, interceptor)
			if err != nil {
				return fmt.Errorf("Failed setting a provider: %w", err)
			}
		} else {
			provisionPodclient, err = publisher.ConnectToProvider(ipfsp2p, deployment, interceptor)
			if err != nil {
				return err
			}
		}

		if uploadSecrets {
			err = publisher.UploadSecrets(cmd.Context(), ipfs, filepath.Dir(podFile), pod, deployment)
			if err != nil {
				return err
			}
		}

		ctrdClient, err := ipcr.GetContainerdClient("k8s.io")
		if err != nil {
			return err
		}

		if uploadImages {
			err = publisher.UploadImages(cmd.Context(), ctrdClient, ipfsApi, pod, deployment)
			if err != nil {
				return err
			}
		}

		err = publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment) // Checkpoint uploads and keys so far
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", err)
		}

		err = publisher.FundPaymentChannel(ethClient, publisherAuth, deployment, fundsInt, unlockTimeInt, debugMintFunds)
		if err != nil {
			return err
		}

		fmt.Printf("PODID is:%v\n", common.BytesToHash(deployment.Payment.PodID))

		response, err := publisher.SendToProvider(cmd.Context(), ipfsp2p, pod, deployment, provisionPodclient)
		if err != nil {
			return err
		}
		// Authorize the application to manipulate the payment channel and fund
		// it to make it able to send transactions
		if authorize {
			err := publisher.AuthorizeAndFundApplication(cmd.Context(), response.(*pb.ProvisionPodResponse), deployment, ethClient, publisherAuth, publisherKey, 1000000000000000000)
			if err != nil {
				return err
			}
		}

		return publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment)
	},
}

var deletePodCmd = &cobra.Command{
	Use:     fmt.Sprintf("delete [%s|deployment.yaml]", publisher.DefaultPodFile),
	Aliases: []string{"undeploy"},
	Short:   "Delete a pod from a local deployment",
	Args:    cobra.MaximumNArgs(1),
	GroupID: "main",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, deploymentFile, _, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		configureDeployment(deployment)

		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connecting to IPFS: %w", err)
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		if publisherKey == "" {
			publisherKey = common.BytesToAddress(deployment.Payment.PublisherAddress).String()
		}

		_, sign, err := ethereum.GetAccountAndSigner(publisherKey, ethClient)
		if err != nil {
			return fmt.Errorf("Could not get ethereum account: %w", err)
		}

		interceptor := pbcon.NewAuthInterceptorClient(deployment, expirationOffset, sign)

		ipfsp2p := tpipfs.NewP2pApi(ipfs, ipfsMultiaddr)

		client, err := publisher.ConnectToProvider(ipfsp2p, deployment, interceptor)
		if err != nil {
			return err
		}

		_, err = publisher.SendToProvider(cmd.Context(), tpipfs.NewP2pApi(ipfs, ipfsMultiaddr), nil, deployment, client)
		if err != nil {
			return err
		}

		return publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment)
	},
}

func init() {
	podCmd.AddCommand(deployPodCmd)
	podCmd.AddCommand(deletePodCmd)

	deployPodCmd.Flags().AddFlagSet(deploymentFlags)
	deployPodCmd.Flags().AddFlagSet(uploadFlags)
	deployPodCmd.Flags().AddFlagSet(fundFlags)
	deployPodCmd.Flags().AddFlagSet(syncFlags)
	deployPodCmd.Flags().AddFlagSet(registryFlags)
	deletePodCmd.Flags().AddFlagSet(deploymentFlags)
	deletePodCmd.Flags().AddFlagSet(syncFlags)

}
