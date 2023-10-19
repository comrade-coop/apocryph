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

var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "Operations related to payment channel",
}

func VerifyContractAddress(address string) error {
	if !slices.Contains(allowedContracts, address) {
		return errors.New(fmt.Sprintf("Contract address (%s) not in the list of allowed contract addresses %v", address, allowedContracts))
	}
	return nil
}

var createChannelCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a payment channel with initial funds",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := common.HexToAddress(PaymentContractAddress)
		unlockTime := big.NewInt(UnlockTime)
		providerAddress := common.HexToAddress(ProviderAddress)
		funds := big.NewInt(Funds)
		tokenAddress := common.HexToAddress(TokenContractAddress)
		podId := common.HexToHash(PodId)

		err := VerifyContractAddress(PaymentContractAddress)
		if err != nil {
			return err
		}

		// derive the Account from the private key
		ks := crypto.CreateKeystore(keystorePath)

		// get an ethclient
		client, err := contracts.Connect(rpc)
		if err != nil {
			return err
		}

		_, clientAuth, err := contracts.DeriveAccountConfigs(ClientKey, Password, exportPassword, client, ks)
		if err != nil {
			return err
		}

		// get a payment contract instance
		payment, err := abi.NewPayment(address, client)
		if err != nil {
			return err
		}

		token, err := abi.NewMockToken(tokenAddress, client)
		if err != nil {
			return err
		}
		// claim token and approve them to the payment contract
		tx, err := token.Mint(clientAuth, funds)
		if err != nil {
			return err
		}

		tx, err = token.Approve(clientAuth, address, funds)
		if err != nil {
			return err
		}

		tx, err = payment.CreateChannel(clientAuth, providerAddress, podId, tokenAddress, unlockTime, funds)
		if err != nil {
			return err
		}

		log.Printf("Payment Channel Created Succefully, tx: %v", tx.Hash())
		return nil
	},
}

func init() {
	paymentCmd.AddCommand(createChannelCmd)

	createChannelCmd.Flags().StringVar(&TokenContractAddress, "tokenAddr", "", "token contract address")
	createChannelCmd.Flags().StringVar(&PodId, "pod-id", "00", "pod id")
	createChannelCmd.Flags().Int64Var(&Funds, "funds", 500, "intial funds")
	createChannelCmd.Flags().Int64Var(&UnlockTime, "unlock-time", 5, "minumum deadline advance duration")
}
