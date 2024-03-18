// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/comrade-coop/apocryph/pkg/abi"
	"github.com/comrade-coop/apocryph/pkg/ethereum"
	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/comrade-coop/apocryph/pkg/publisher"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Operations related to registry",
}

var getTablesCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Pricing tables filtered by the provided prices and provider regions",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, _, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connecting to IPFS: %w", err)
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		registryContract := common.HexToAddress(registryContractAddress)
		var tokenContract common.Address
		if tokenContractAddress != "" {
			tokenContract = common.HexToAddress(tokenContractAddress)
		} else if paymentContractAddress != "" {
			paymentContract := common.HexToAddress(paymentContractAddress)
			payment, err := abi.NewPayment(paymentContract, ethClient)
			if err != nil {
				return fmt.Errorf("Failed instantiating payment contract: %w", err)
			}
			tokenContract, err = payment.Token(&bind.CallOpts{})
			if err != nil {
				return fmt.Errorf("Failed getting payment contract token: %w", err)
			}
		} else {
			return fmt.Errorf("Either a token contract or a payment contract must be set")
		}

		tables, err := publisher.GetPricingTables(ethClient, registryContract, tokenContract)
		if err != nil {
			return err
		}

		filter, err := getRegistryTableFilter()
		if err != nil {
			return err
		}

		filteredTables := publisher.FilterPricingTables(tables, filter)

		if len(filteredTables) == 0 {
			return fmt.Errorf("no table found by given filter")
		}

		PrintTableInfo(filteredTables, ethClient, registryContract)

		allProviders, err := publisher.GetProviderHostInfos(ipfs, ethClient, registryContract, filteredTables)
		if err != nil {
			return err
		}
		filteredInfos, err := publisher.FilterProviderHostInfos(region, "", allProviders)
		if err != nil {
			return err
		}
		if len(filteredInfos) == 0 {
			return fmt.Errorf("could not find providers with the given filter")
		}

		PrintProvidersInfo(filteredInfos)

		return nil
	},
}

func PrintProvidersInfo(filteredInfos publisher.ProviderHostInfoList) {
	infoTable := tablewriter.NewWriter(os.Stdout)
	infoTable.SetHeader([]string{"Id", "Regions", "Addresses"})
	for id, info := range filteredInfos {
		row := []string{id.Hex(), formatRegions(info.Regions), strings.Join(info.Multiaddrs, "\n")}
		infoTable.Append(row)
	}
	infoTable.SetRowLine(true)
	fmt.Printf("Providers:\n")
	infoTable.Render()
}

func PrintTableInfo(filteredTables publisher.PricingTableList, ethClient *ethclient.Client, registryContract common.Address) error {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Id", "CPU", "RAM", "STORAGE", "BANDWIDTH EGRESS", "BANDWIDTH INGRESS", "CPU MODEL", "TEE TECHNOLOGY", "Providers"})
	for _, t := range filteredTables {
		subscribers, err := publisher.GetPricingTableSubscribers(ethClient, registryContract, t.Id)
		if err != nil {
			return err
		}
		// if table has no subscribers skip printing it
		if len(subscribers) == 0 {
			continue
		}
		var subs string = ""
		for _, subscriber := range subscribers {
			subs = subs + subscriber.String() + "\n"
		}
		row := []string{t.Id.String(), t.CpuPrice.String(), t.RamPrice.String(),
			t.StoragePrice.String(), t.BandwidthEgressPrice.String(),
			t.BandwidthIngressPrice.String(), t.Cpumodel, t.TeeType, subs}
		table.Append(row)
	}
	table.Render()
	return nil
}

func formatRegions(regions []*pb.HostInfo_Region) string {
	result := ""
	for _, region := range regions {
		result += fmt.Sprintf("%s-%s-%d\n", region.Name, region.Zone, region.Num)
	}
	return result
}

func init() {
	registryCmd.AddCommand(getTablesCmd)
	getTablesCmd.Flags().AddFlagSet(registryFlags)
}
