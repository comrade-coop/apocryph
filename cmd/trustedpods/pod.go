package main

import (
	"crypto/sha256"
	"fmt"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/provider"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Operations related to with raw pod manifests",
}

// if no provider is selected, Fetches providers based on registry args
func fetchAndFilterProviders() ([]*provider.HostInfo, error) {

	if providerPeer != "" {
		return []*provider.HostInfo{{Id: providerPeer}}, nil
	}

	var availableProviders []*provider.HostInfo
	// Get pricing table filtered by registryFlags
	tables, registry, ipfsApi, ethClient, err := publisher.GetRegistryComponents(ipfsApi, ethereumRpc, registryContractAddress, tokenContractAddress)
	if err != nil {
		return nil, err
	}

	//filter tables
	filteredTables := publisher.FilterTables(tables, getFilter())

	// Get available providers
	if len(filteredTables) == 0 {
		return nil, fmt.Errorf("no table found by filter")
	}

	availableProviders, err = publisher.GetProvidersHostingInfo(ipfsApi, ethClient, registry, filteredTables)
	if err != nil {
		return nil, err
	}

	// filter Providers
	availableProviders, err = publisher.FilterProviders(region, providerPeer, availableProviders)
	if err != nil {
		return nil, err
	}

	return availableProviders, nil
}

func configureDeployment(deployment *pb.Deployment) error {
	if providerEthAddress != "" && providerPeer != "" {
		deployment.Provider = &pb.ProviderConfig{
			EthereumAddress: common.HexToAddress(providerEthAddress).Bytes(),
			Libp2PAddress:   providerPeer,
		}
	}
	if deployment.Payment == nil {
		deployment.Payment = &pb.PaymentChannelConfig{}
	}
	if paymentContractAddress != "" {
		deployment.Payment.PaymentContractAddress = common.HexToAddress(paymentContractAddress).Bytes()
	}
	if podId != "" {
		deployment.Payment.PodID = common.HexToHash(podId).Bytes()
	}
	if deployment.Payment.PodID == nil {
		podFileNameHash := sha256.Sum256([]byte(deployment.PodManifestFile))
		deployment.Payment.PodID = podFileNameHash[:]
	}
	return nil
}

func init() {
	rootCmd.AddCommand(podCmd)
	rootCmd.AddCommand(registryCmd)

	podCmd.AddGroup(&cobra.Group{
		ID:    "main",
		Title: "Commands",
	})
	podCmd.AddGroup(&cobra.Group{
		ID:    "lowlevel",
		Title: "Low-level commands",
	})

	podCmd.PersistentFlags().AddFlagSet(podFlags)
}
