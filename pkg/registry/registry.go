// SPDX-License-Identifier: GPL-3.0

package registry

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/proto/protoconnect"
	"github.com/comrade-coop/trusted-pods/pkg/provider"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"gopkg.in/yaml.v2"
)

type PricingTable struct {
	Token                 string
	Id                    big.Int
	CpuPrice              big.Int
	RamPrice              big.Int
	StoragePrice          big.Int
	BandwidthEgressPrice  big.Int
	BandwidthIngressPrice big.Int
	CpuModel              string
	TeeType               string
}

type Filter struct {
	CpuPrice *big.Int
}

func GetPricingTables(ethClient *ethclient.Client, registry *abi.Registry, tokenAddress string) ([]*abi.RegistryNewPricingTable, error) {
	iter, err := registry.FilterNewPricingTable(&bind.FilterOpts{Start: 0}, []common.Address{common.HexToAddress(tokenAddress)}, nil)
	if err != nil {
		return nil, err
	}
	var tables []*abi.RegistryNewPricingTable
	for iter.Next() {
		if iter.Error() != nil {
			return nil, iter.Error()
		}
		tables = append(tables, iter.Event)
	}

	return tables, nil
}

func FilterTables(tables []*abi.RegistryNewPricingTable, filter *abi.RegistryNewPricingTable) []*abi.RegistryNewPricingTable {
	var filteredTables []*abi.RegistryNewPricingTable

	for _, table := range tables {
		if match(table, filter) {
			filteredTables = append(filteredTables, table)
		}
	}
	return filteredTables
}

func match(table, filter *abi.RegistryNewPricingTable) bool {
	// Check each field in the filter. If the field is nil, consider it a match. Otherwise, compare the values.
	if filter.Token.Cmp(table.Token) != 0 {
		return false
	}
	if filter.Id != nil && filter.Id.Cmp(table.Id) != 0 {
		return false
	}
	if filter.CpuPrice != nil && filter.CpuPrice.Cmp(table.CpuPrice) != 0 {
		return false
	}
	if filter.RamPrice != nil && filter.RamPrice.Cmp(table.RamPrice) != 0 {
		return false
	}
	if filter.StoragePrice != nil && filter.StoragePrice.Cmp(table.StoragePrice) != 0 {
		return false
	}
	if filter.BandwidthEgressPrice != nil && filter.BandwidthEgressPrice.Cmp(table.BandwidthEgressPrice) != 0 {
		return false
	}
	if filter.BandwidthIngressPrice != nil && filter.BandwidthIngressPrice.Cmp(table.BandwidthIngressPrice) != 0 {
		return false
	}
	if filter.Cpumodel != "" && !strings.Contains(table.Cpumodel, filter.Cpumodel) {
		return false
	}
	if filter.TeeType != "" && !strings.Contains(table.TeeType, filter.TeeType) {
		return false
	}

	return true
}

func GetTableSubscribers(ethClient *ethclient.Client, registry *abi.Registry, tableId *big.Int) ([]common.Address, error) {
	iter, err := registry.FilterSubscribed(&bind.FilterOpts{Start: 0}, []*big.Int{tableId}, nil)
	if err != nil {
		return nil, err
	}
	var subscribers []common.Address
	for iter.Next() {
		if iter.Error() != nil {
			return nil, iter.Error()
		}
		subscribed, err := registry.RegistryCaller.IsSubscribed(&bind.CallOpts{Pending: false}, iter.Event.Provider, tableId)
		if err != nil {
			return nil, err
		}
		if subscribed {
			subscribers = append(subscribers, iter.Event.Provider)
		}
	}
	return subscribers, nil
}

func GetProviderHostInfo(ipfs *rpc.HttpApi, registry *abi.Registry, p common.Address) (*provider.HostInfo, error) {
	id, err := registry.Providers(&bind.CallOpts{}, p)
	if err != nil {
		return nil, err
	}

	infoCid, err := cid.Decode(id)
	if err != nil {
		return nil, fmt.Errorf("failed converting cid to bytes: %v", err)
	}

	infoBytes, err := tpipfs.GetBytes(ipfs, infoCid)
	if err != nil {
		return nil, fmt.Errorf("failed retreiving provider info from ipfs: %v", err)
	}

	providerInfo := &provider.HostInfo{}

	err = yaml.Unmarshal(infoBytes, providerInfo)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling provider Info : %v", err)
	}

	return providerInfo, nil
}

func FilterProviders(region, peerId string, providers []*provider.HostInfo) ([]*provider.HostInfo, error) {
	if region == "" {
		return providers, nil
	}
	regionParts := strings.Split(region, "-")
	if len(regionParts) != 3 {
		return nil, fmt.Errorf("Malformed filter: region")
	}
	parsedRegion := provider.Region{}
	parsedRegion.Name = regionParts[0]
	parsedRegion.Zone = regionParts[1]
	num, err := strconv.ParseUint(regionParts[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Could not convert region number: %v", err)
	}
	parsedRegion.Num = uint32(num)
	var filteredProviders []*provider.HostInfo
	for _, provider := range providers {
		for _, region := range provider.Regions {
			if region.IsEqual(parsedRegion) {
				if peerId != "" && peerId != provider.Id {
					continue
				}
				filteredProviders = append(filteredProviders, provider)
			}
		}
	}
	return filteredProviders, nil
}

func GetRegistryComponents(ipfsApi, ethereumRpc, registryContractAddress, tokenContractAddress string) ([]*abi.RegistryNewPricingTable, *abi.Registry, *rpc.HttpApi, *ethclient.Client, error) {
	ipfs, _, err := tpipfs.GetIpfsClient(ipfsApi)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Failed connecting to IPFS: %w", err)
	}

	ethClient, err := ethereum.GetClient(ethereumRpc)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	registry, err := abi.NewRegistry(common.HexToAddress(registryContractAddress), ethClient)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tables, err := GetPricingTables(ethClient, registry, tokenContractAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if len(tables) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("no pricing table found for token %v", tokenContractAddress)
	}

	return tables, registry, ipfs, ethClient, nil
}

func GetProvidersHostingInfo(ipfs *rpc.HttpApi, ethClient *ethclient.Client, registry *abi.Registry, tables []*abi.RegistryNewPricingTable) ([]*provider.HostInfo, error) {
	var hostingInfoTables []*provider.HostInfo
	for _, table := range tables {
		subscribers, err := GetTableSubscribers(ethClient, registry, table.Id)
		if err != nil {
			return nil, err
		}
		for _, subscriber := range subscribers {
			info, err := GetProviderHostInfo(ipfs, registry, subscriber)
			if err != nil {
				return nil, err
			}
			hostingInfoTables = append(hostingInfoTables, info)
		}
	}
	return hostingInfoTables, nil
}

func SetFirstConnectingProvider(ipfsp2p *tpipfs.P2pApi, availableProviders []*provider.HostInfo, deployment *pb.Deployment, interceptor *protoconnect.AuthInterceptorClient) (*publisher.P2pProvisionPodServiceClient, error) {
	deployment.Provider = &pb.ProviderConfig{}
	for _, provider := range availableProviders {
		deployment.Provider.Libp2PAddress = provider.MultiAddresses[0].Value
		deployment.Provider.EthereumAddress = common.HexToAddress(provider.Id).Bytes()
		client, err := publisher.ConnectToProvider(ipfsp2p, deployment, interceptor)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		return client, nil
	}
	return nil, fmt.Errorf("No provider is reachable")
}
