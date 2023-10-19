package main

import (
	"errors"
	"fmt"
	"log"
	"math/big"

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
		minAdvanceDuration := big.NewInt(MinAdvanceDuration)
		providerAddress := common.HexToAddress(ProviderAddress)
		funds := big.NewInt(Funds)
		price := big.NewInt(PricePerExecution)
		deadline := big.NewInt(Deadline)
		tokenAddress := common.HexToAddress(TokenContractAddress)

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
		payment, err := contracts.GetContractInstance(client, address)
		if err != nil {
			return err
		}

		token, err := contracts.GetTokenContractInstance(client, tokenAddress)
		if err != nil {
			return err
		}

		// claim token and approve them to the payment contract
		tx, err := token.ClaimTokens(clientAuth, funds)
		if err != nil {
			return err
		}
		tx, err = token.Approve(clientAuth, address, funds)
		if err != nil {
			return err
		}

		tx, err = contracts.CreatePaymentChannel(clientAuth, payment, providerAddress, tokenAddress, funds, deadline, minAdvanceDuration, price)
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
	createChannelCmd.Flags().Int64Var(&Funds, "funds", 500, "intial funds")
	createChannelCmd.Flags().Int64Var(&Deadline, "deadline", 3275538098, "Deadline for payment channel expiration date")
	createChannelCmd.Flags().Int64Var(&PricePerExecution, "execPrice", 5, "Price per execution")
	createChannelCmd.Flags().Int64Var(&MinAdvanceDuration, "minDuration", 5, "minumum deadline advance duration")
}
