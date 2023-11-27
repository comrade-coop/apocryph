package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Operations related to registry",
}

func getFilter() *abi.RegistryNewPricingTable {
	var id, cPrice, rPrice, sPrice, bePrice, biPrice *big.Int
	cPrice, _ = (&big.Int{}).SetString(cpuPrice, 10)
	rPrice, _ = (&big.Int{}).SetString(ramPrice, 10)
	sPrice, _ = (&big.Int{}).SetString(storagePrice, 10)
	bePrice, _ = (&big.Int{}).SetString(bandwidthEPrice, 10)
	biPrice, _ = (&big.Int{}).SetString(bandwidthInPrice, 10)
	id, _ = (&big.Int{}).SetString(tableId, 10)

	filter := abi.RegistryNewPricingTable{Token: common.HexToAddress(tokenContractAddress), Id: id,
		CpuPrice:              cPrice,
		RamPrice:              rPrice,
		StoragePrice:          sPrice,
		BandwidthEgressPrice:  bePrice,
		BandwidthIngressPrice: biPrice,
		Cpumodel:              cpuModel,
		TeeType:               teeType}
	return &filter
}

var getTablesCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Pricing tables filtered by the provided prices and provider regions",
	RunE: func(cmd *cobra.Command, args []string) error {

		tables, registry, ipfs, ethClient, err := publisher.GetRegistryComponents(ipfsApi, ethereumRpc, registryContractAddress, tokenContractAddress)
		if err != nil {
			return err
		}

		filteredTables := publisher.FilterTables(tables, getFilter())
		table := tablewriter.NewWriter(os.Stdout)
		var allSubscribers []common.Address

		table.SetHeader([]string{"Id", "CPU", "RAM", "STORAGE", "BANDWIDTH EGRESS", "BANDWIDTH INGRESS", "CPU MODEL", "TEE TECHNOLOGY", "Providers"})
		for _, t := range filteredTables {
			subscribers, err := publisher.GetTableSubscribers(ethClient, registry, t.Id)
			if err != nil {
				return err
			}
			var subs string = ""
			for _, subscriber := range subscribers {
				subs = subs + subscriber.String() + "\n"
				if !slices.Contains(allSubscribers, subscriber) {
					allSubscribers = append(allSubscribers, subscriber)
				}
			}
			// if table has no subscribers skip printing it
			if len(subscribers) != 0 {
				row := []string{t.Id.String(), t.CpuPrice.String(), t.RamPrice.String(),
					t.StoragePrice.String(), t.BandwidthEgressPrice.String(),
					t.BandwidthIngressPrice.String(), t.Cpumodel, t.TeeType, subs}
				table.Append(row)
			}
		}
		if len(filteredTables) == 0 {
			return fmt.Errorf("no table found by given filter")
		}
		table.Render()

		infoTable := tablewriter.NewWriter(os.Stdout)
		infoTable.SetHeader([]string{"Id", "Regions", "Addresses"})
		allProviders, err := publisher.GetProvidersHostingInfo(ipfs, ethClient, registry, filteredTables)
		if err != nil {
			return err
		}
		filteredInfos, err := publisher.FilterProviders(region, "", allProviders)
		if err != nil {
			return err
		}
		if len(filteredInfos) == 0 {
			return fmt.Errorf("could not find providers with the given filter")
		}
		for _, info := range filteredInfos {
			row := []string{info.Id, info.FormatRegions(), info.FormatAddresses()}
			infoTable.Append(row)
		}
		infoTable.SetRowLine(true)
		fmt.Printf("\nProviders:\n")
		infoTable.Render()

		return nil
	},
}

func init() {
	registryCmd.AddCommand(getTablesCmd)
	getTablesCmd.Flags().AddFlagSet(registryFlags)
}
