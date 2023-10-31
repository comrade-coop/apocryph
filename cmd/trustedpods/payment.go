package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var ethereumRpc string
var publisherKey string
var paymentContractAddress string
var providerEthAddress string
var podId string
var tokenContractAddress string
var unlockTime int64
var funds string

var paymentCmd = &cobra.Command{
	Use:   "payment",
	Short: "Operations related to payment channel",
}

func createPaymentChannel(podIdBytes common.Hash) (*pb.PaymentChannel, error) {
	paymentContract := common.HexToAddress(paymentContractAddress)
	unlockTimeInt := big.NewInt(unlockTime)
	provider := common.HexToAddress(providerEthAddress)
	tokenContract := common.HexToAddress(tokenContractAddress)

	fundsInt, _ := (&big.Int{}).SetString(funds, 10)
	if fundsInt == nil {
		return nil, fmt.Errorf("Invalid number passed for funds: %s", funds)
	}

	client, err := ethereum.GetClient(ethereumRpc)
	if err != nil {
		return nil, err
	}

	publisherAuth, err := ethereum.GetAccount(publisherKey, client)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	// get a payment contract instance
	payment, err := abi.NewPayment(paymentContract, client)
	if err != nil {
		return nil, err
	}

	token, err := abi.NewIERC20(tokenContract, client)
	if err != nil {
		return nil, err
	}

	if fundsInt.Cmp(common.Big0) != 0 {
		tx, err := token.Approve(publisherAuth, paymentContract, fundsInt)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Token approval successful! %v\n", tx.Hash())

		tx, err = payment.CreateChannel(publisherAuth, provider, podIdBytes, tokenContract, unlockTimeInt, fundsInt)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Payment channel created! %v\n", tx.Hash())
	}

	return &pb.PaymentChannel{
		ChainID:          chainID.Bytes(),
		ContractAddress:  paymentContract.Bytes(),
		PublisherAddress: publisherAuth.From.Bytes(),
		ProviderAddress:  provider.Bytes(),
		PodID:            podIdBytes.Bytes(),
		TokenAddress:     tokenContract.Bytes(),
	}, nil
}

var createPaymentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a payment channel with initial funds",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := createPaymentChannel(common.HexToHash(podId))
		return err
	},
}

var mintPaymentCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint some amount of tokens (debug)",
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenContract := common.HexToAddress(tokenContractAddress)
		fundsInt, _ := (&big.Int{}).SetString(funds, 10)
		if fundsInt == nil {
			return fmt.Errorf("Invalid number passed for funds: %s", funds)
		}

		// get an ethclient
		client, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		// derive the Account from the private key
		publisherAuth, err := ethereum.GetAccount(publisherKey, client)
		if err != nil {
			return err
		}

		token, err := abi.NewMockToken(tokenContract, client)
		if err != nil {
			return err
		}

		// claim token and approve them to the payment contract
		tx, err := token.Mint(publisherAuth, fundsInt)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Mint successful! %v", tx.Hash())
		return nil
	},
}

func init() {
	paymentCmd.AddCommand(createPaymentCmd)
	paymentCmd.AddCommand(mintPaymentCmd)

	paymentCmd.PersistentFlags().StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "ethereum rpc node")
	paymentCmd.PersistentFlags().StringVar(&publisherKey, "ethereum-key", "", "account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")
	paymentCmd.PersistentFlags().StringVar(&tokenContractAddress, "token", "", "token contract address")
	createPaymentCmd.Flags().StringVar(&paymentContractAddress, "payment-contract", "", "payment contract address")
	createPaymentCmd.Flags().StringVar(&providerEthAddress, "provider-eth", "", "provider public address")
	createPaymentCmd.Flags().StringVar(&podId, "pod-id", "00", "pod id")
	createPaymentCmd.Flags().StringVar(&funds, "funds", "5000000000000000000", "intial funds")
	createPaymentCmd.Flags().Int64Var(&unlockTime, "unlock-time", 5*60, "time for unlocking tokens (in seconds)")
	mintPaymentCmd.Flags().StringVar(&funds, "funds", "5000000000000000000", "amount to mint")
}
