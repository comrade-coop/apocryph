package main

import (
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
var podID string
var checkContractCmd = &cobra.Command{
	Use:   "check",
	Short: "Check whether a payment contract has a permissible hash",
	RunE: func(cmd *cobra.Command, args []string) error {
		address := common.HexToAddress(PaymentContractAddress)
		if !slices.Contains(allowedContractCodeHashes, address.Hex()) {
			return errors.New(fmt.Sprintf("Contract code hash (%s) not in the list of allowed code hashes %v", address, allowedContractCodeHashes))
		}

		fmt.Println("Correct Contract")
		return nil
	},
}

var claimCmd = &cobra.Command{
	Use:   "claim <amount>",
	Short: "Claim funds from a payment contract",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		publisherAddress := common.HexToAddress(PublisherAddress)
		tokenAddress := common.HexToAddress(TokenContractAddress)
		paymentAddress := common.HexToAddress(PaymentContractAddress)
		podIDBytes := common.HexToHash(podID)
		amount := big.NewInt(10)
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
		payment, err := abi.NewPayment(paymentAddress, client)

		_, err = payment.Withdraw(providerAuth, publisherAddress, podIDBytes, tokenAddress, amount, common.Address{})
		if err != nil {
			return err
		}
		log.Println("Amount transfered successfully")
		return nil
	},
}

func init() {
	contractCmd.AddCommand(checkContractCmd)
	contractCmd.AddCommand(claimCmd)

	claimCmd.Flags().StringVar(&PublisherAddress, "cAddr", "", "client public address")
	claimCmd.Flags().StringVar(&TokenContractAddress, "tokenAddr", "", "token contract address")
	claimCmd.Flags().StringVar(&podID, "id", "00", "pod ID (hex)")

}
