// SPDX-License-Identifier: GPL-3.0

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/comrade-coop/apocryph/pkg/abi"
	"github.com/comrade-coop/apocryph/pkg/ethereum"
	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/comrade-coop/apocryph/pkg/resource"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var registryContractAddress string
var tokenContractAddress string
var hostInfoContents string
var hostInfoFormat string
var cpuModel string
var teeType string

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Operations related to registry",
}

var registerSelfCmd = &cobra.Command{
	Use:   "self",
	Short: "upload provider info to ipfs and register in the registry smart contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, _, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		providerAuth, _, err := ethereum.GetAccountAndSigner(providerKey, ethClient)
		if err != nil {
			return err
		}

		registry, err := abi.NewRegistry(common.HexToAddress(registryContractAddress), ethClient)
		if err != nil {
			return err
		}

		hostInfo, err := getHostInfo(cmd.Context(), ipfs)
		if err != nil {
			return err
		}

		cid, err := tpipfs.AddProtobufFile(ipfs, hostInfo)
		if err != nil {
			return err
		}

		tx, err := registry.RegisterProvider(providerAuth, cid.String())
		if err != nil {
			return err
		}

		log.Println("Provider Info Cid:\n", cid.String())
		log.Printf("Registered provider, TX Hash:%v\n", tx.Hash())

		return nil
	},
}

var unsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe <TABLE_ID>",
	Short: "unsubscribe from pricing table",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		providerAuth, _, err := ethereum.GetAccountAndSigner(providerKey, ethClient)
		if err != nil {
			return err
		}

		registry, err := abi.NewRegistry(common.HexToAddress(registryContractAddress), ethClient)
		if err != nil {
			return err
		}

		tableId, _ := (&big.Int{}).SetString(id, 10)
		_, err = registry.Unsubscribe(providerAuth, tableId)
		if err != nil {
			return err
		}

		return nil
	},
}
var subscribeCmd = &cobra.Command{
	Use:   "subscribe <TABLE_ID>",
	Short: "subscribe provider in pricing table (Requires provider to be registered)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		providerAuth, _, err := ethereum.GetAccountAndSigner(providerKey, ethClient)
		if err != nil {
			return err
		}

		registry, err := abi.NewRegistry(common.HexToAddress(registryContractAddress), ethClient)
		if err != nil {
			return err
		}

		tableId, _ := (&big.Int{}).SetString(id, 10)
		_, err = registry.Subscribe(providerAuth, tableId)
		if err != nil {
			return err
		}

		return nil
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "uploads provider info to ipfs and register in the registry contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipfs, _, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return err
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		providerAuth, _, err := ethereum.GetAccountAndSigner(providerKey, ethClient)
		if err != nil {
			return err
		}

		pricingTables, err := openPricingTables()
		if err != nil {
			return err
		}

		registry, err := abi.NewRegistry(common.HexToAddress(registryContractAddress), ethClient)
		if err != nil {
			return err
		}

		hostInfo, err := getHostInfo(cmd.Context(), ipfs)
		if err != nil {
			return err
		}

		cid, err := tpipfs.AddProtobufFile(ipfs, hostInfo)
		if err != nil {
			return err
		}

		tx, err := registry.RegisterProvider(providerAuth, cid.String())
		if err != nil {
			return fmt.Errorf("Could not register Provider: %v", err)
		}

		prices, err := resource.GetTablesPrices(pricingTables)
		if err != nil {
			return err
		}
		for i := 0; i < len(prices); i++ {
			// register pricing tables in the contract
			tx, err = registry.RegisterPricingTable(providerAuth, common.HexToAddress(tokenContractAddress), prices[i], cpuModel, teeType)
			if err != nil {
				return err
			}
			fmt.Printf("table registered, TX Hash %v\n", tx.Hash())
		}

		return nil
	},
}

func getHostInfo(ctx context.Context, ipfs *rpc.HttpApi) (*pb.HostInfo, error) {
	if hostInfoContents == "" {
		return nil, fmt.Errorf("Empty host info")
	}
	hostInfo := &pb.HostInfo{}
	err := pb.Unmarshal(hostInfoFormat, []byte(hostInfoContents), hostInfo)
	if err != nil {
		return nil, err
	}
	key, err := ipfs.Key().Self(ctx)
	if err != nil {
		return nil, err
	}
	peerId := key.ID()
	hostInfo.Multiaddrs = append(hostInfo.Multiaddrs, fmt.Sprintf("/p2p/%s", peerId.String()))
	return hostInfo, nil
}

func init() {
	registryFlags := &pflag.FlagSet{}

	registryFlags.StringVar(&ipfsApi, "ipfs", "", "multiaddr where the ipfs/kubo api can be accessed (leave blank to use the daemon running in IPFS_PATH)")
	registryFlags.StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "client public address")
	registryFlags.StringVar(&registryContractAddress, "registry-contract", "", "registry contract address")
	registryFlags.StringVar(&tokenContractAddress, "token-contract", "", "token contract address")
	registryFlags.StringVar(&providerKey, "ethereum-key", "", "provider account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")

	registerSelfCmd.Flags().AddFlagSet(registryFlags)
	registerCmd.Flags().AddFlagSet(registryFlags)
	subscribeCmd.Flags().AddFlagSet(registryFlags)
	unsubscribeCmd.Flags().AddFlagSet(registryFlags)

	AddConfig("info.contents", &hostInfoContents, "", "host info section from config file")
	AddConfig("info.format", &hostInfoFormat, "", "host info format")
	AddConfig("cpu_model", &cpuModel, "", "provider CPU model ")
	AddConfig("tee_type", &teeType, "", "type of tee technology that is used (Secure Enclaes/CVMs/...etc)")

	registryCmd.AddCommand(registerCmd)
	registryCmd.AddCommand(registerSelfCmd)
	registryCmd.AddCommand(subscribeCmd)
	registryCmd.AddCommand(unsubscribeCmd)
}
