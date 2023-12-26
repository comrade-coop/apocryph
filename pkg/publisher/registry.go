// SPDX-License-Identifier: GPL-3.0

package publisher

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/comrade-coop/apocryph/pkg/abi"
	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/comrade-coop/apocryph/pkg/proto/protoconnect"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/protobuf/proto"
)

type PricingTableList []*abi.RegistryNewPricingTable

func GetPricingTables(ethClient *ethclient.Client, registryContract common.Address, tokenAddress common.Address) (PricingTableList, error) {
	registry, err := abi.NewRegistry(registryContract, ethClient)
	if err != nil {
		return nil, err
	}

	iter, err := registry.FilterNewPricingTable(&bind.FilterOpts{Start: 0}, []common.Address{tokenAddress}, nil)
	if err != nil {
		return nil, err
	}
	var tables PricingTableList
	for iter.Next() {
		if iter.Error() != nil {
			return nil, iter.Error()
		}
		tables = append(tables, iter.Event)
	}

	return tables, nil
}

func FilterPricingTables(tables PricingTableList, filter *abi.RegistryNewPricingTable) PricingTableList {
	var filteredTables PricingTableList

	for _, table := range tables {
		if matchPricingTable(table, filter) {
			filteredTables = append(filteredTables, table)
		}
	}
	return filteredTables
}

func matchPricingTable(table, filter *abi.RegistryNewPricingTable) bool {
	// Check each field in the filter. If the field is nil, consider it a match. Otherwise, compare the values.
	if filter.Token != common.BytesToAddress(nil) && filter.Token.Cmp(table.Token) != 0 {
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

func GetPricingTableSubscribers(ethClient *ethclient.Client, registryContract common.Address, tableId *big.Int) ([]common.Address, error) {
	registry, err := abi.NewRegistry(registryContract, ethClient)
	if err != nil {
		return nil, err
	}

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

func GetProviderHostInfo(ipfs iface.CoreAPI, ethClient *ethclient.Client, registryContract common.Address, providerId common.Address) (*pb.HostInfo, error) {
	registry, err := abi.NewRegistry(registryContract, ethClient)
	if err != nil {
		return nil, err
	}

	id, err := registry.Providers(&bind.CallOpts{}, providerId)
	if err != nil {
		return nil, err
	}

	infoCid, err := cid.Decode(id)
	if err != nil {
		return nil, fmt.Errorf("failed converting cid to bytes: %w", err)
	}

	providerInfo := &pb.HostInfo{}

	err = tpipfs.GetProtobufFile(ipfs, infoCid, providerInfo)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving provider info from ipfs: %w", err)
	}

	return providerInfo, nil
}

type ProviderHostInfoList map[common.Address]*pb.HostInfo

func GetProviderHostInfos(ipfs iface.CoreAPI, ethClient *ethclient.Client, registryContract common.Address, tables PricingTableList) (ProviderHostInfoList, error) {
	providers := make(ProviderHostInfoList)
	for _, table := range tables {
		subscribers, err := GetPricingTableSubscribers(ethClient, registryContract, table.Id)
		if err != nil {
			return nil, err
		}
		for _, subscriber := range subscribers {
			info, err := GetProviderHostInfo(ipfs, ethClient, registryContract, subscriber)
			if err != nil {
				return nil, err
			}
			providers[subscriber] = info
		}
	}
	return providers, nil
}

func FilterProviderHostInfos(region string, peerId string, providers ProviderHostInfoList) (ProviderHostInfoList, error) {
	if region == "" {
		return providers, nil
	}
	regionParts := strings.Split(region, "-")
	if len(regionParts) != 3 {
		return nil, fmt.Errorf("Malformed filter: region (given %v, expected a string like XX-YY-1)", region)
	}
	parsedRegion := &pb.HostInfo_Region{}
	parsedRegion.Name = regionParts[0]
	parsedRegion.Zone = regionParts[1]
	num, err := strconv.ParseUint(regionParts[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Could not convert region number: %w", err)
	}
	parsedRegion.Num = uint32(num)

	filteredProviders := make(ProviderHostInfoList)
	if peerId != "" {
		id := common.HexToAddress(peerId)
		if id.Hex() != peerId {
			return nil, fmt.Errorf("Malformed filter: peerId (given: %v, expected Ethereum address)", peerId)
		}
		if provider, ok := providers[id]; ok {
			for _, region := range provider.Regions {
				if proto.Equal(region, parsedRegion) {
					filteredProviders[id] = provider
					break
				}
			}
		}
	} else {
		for id, provider := range providers {
			for _, region := range provider.Regions {
				if proto.Equal(region, parsedRegion) {
					filteredProviders[id] = provider
					break
				}
			}
		}
	}
	return filteredProviders, nil
}

func SetFirstConnectingProvider(ipfsp2p *tpipfs.P2pApi, availableProviders ProviderHostInfoList, deployment *pb.Deployment, interceptor *protoconnect.AuthInterceptorClient) (*P2pProvisionPodServiceClient, error) {
	deployment.Provider = &pb.ProviderConfig{}
	for address, provider := range availableProviders {
		found := false
		for _, addr := range provider.Multiaddrs {
			maddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				fmt.Printf("warn: %v\n", err)
				continue
			}
			_, last := multiaddr.SplitLast(maddr)
			if last.Protocol().Name == "p2p" {
				found = true
				deployment.Provider.Libp2PAddress = last.Value()
				break
			}
		}
		if !found {
			err := fmt.Errorf("None of the supplied provider multiaddresses is a /p2p/.. address")
			fmt.Printf("warn: %v\n", err)
			continue
		}
		deployment.Provider.EthereumAddress = address.Bytes()
		client, err := ConnectToProvider(ipfsp2p, deployment, interceptor)
		if err != nil {
			fmt.Printf("warn: %v\n", err)
			continue
		}
		return client, nil
	}
	return nil, fmt.Errorf("No provider is reachable")
}
