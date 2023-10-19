package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/comrade-coop/trusted-pods/pkg/contracts"
	"github.com/comrade-coop/trusted-pods/pkg/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"k8s.io/utils/strings/slices"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Operations related to managing contracts",
}

var ProviderKey string
var chainRpc string
var podID int64
var checkContractCmd = &cobra.Command{
	Use:   "check <contract>",
	Short: "Check whether a payment contract has a permissible hash",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]
		if !slices.Contains(allowedContractCodeHashes, address) {
			return errors.New(fmt.Sprintf("Contract code hash (%s) not in the list of allowed code hashes %v", address, allowedContractCodeHashes))
		}

		// get an ethclient
		client, err := contracts.Connect(chainRpc)
		if err != nil {
			return err
		}

		result, err := client.CodeAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			return err
		}

		byteCode := hex.EncodeToString(result)
		// I assume this is irrelevant and its just compile time related metadata
		currentByteCode := abi.PaymentMetaData.Bin[60:]
		if currentByteCode == byteCode {
			fmt.Println("Correct Contract")
			return nil
		}
		return errors.New("Current Contract bytecode does not match deployed contract bytecode")
	},
}

var claimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Claim funds from a payment contract",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		clientAddress := common.HexToAddress(ClientAddress)
		tokenAddress := common.HexToAddress(TokenContractAddress)
		paymentAddress := common.HexToAddress(PaymentContractAddress)
		// create or get a keystore
		ks := crypto.CreateKeystore(keystorePath)

		// get an ethclient
		client, err := contracts.Connect(chainRpc)
		if err != nil {
			return err
		}
		_, providerAuth, err := contracts.DeriveAccountConfigs(ProviderKey, Password, exportPassword, client, ks)
		if err != nil {
			return err
		}

		// get a contract instance
		payment, err := contracts.GetContractInstance(client, paymentAddress)

		_, err = contracts.Withdraw(providerAuth, payment, clientAddress, big.NewInt(podID), tokenAddress)
		if err != nil {
			return err
		}
		log.Println("Owed Amount transfered successfully")
		return nil
	},
}

func init() {
	contractCmd.AddCommand(checkContractCmd)
	contractCmd.AddCommand(claimCmd)

	claimCmd.Flags().StringVar(&ClientAddress, "cAddr", "", "client public address")
	claimCmd.Flags().StringVar(&TokenContractAddress, "tokenAddr", "", "token contract address")
	claimCmd.Flags().Int64Var(&podID, "id", 1, "pod ID")

}
