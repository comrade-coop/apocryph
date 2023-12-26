// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/resource"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Operations related to managing contracts",
}

var ethereumRpc string
var providerKey string
var paymentContractAddress string
var publisherEthAddress string
var podId string
var metricsTotal int64

func getPaymentChannelProto(providerAuth *bind.TransactOpts, chainID *big.Int) *pb.PaymentChannel {
	return &pb.PaymentChannel{
		ChainID:          chainID.Bytes(),
		ContractAddress:  common.HexToAddress(paymentContractAddress).Bytes(),
		PublisherAddress: common.HexToAddress(publisherEthAddress).Bytes(),
		ProviderAddress:  providerAuth.From.Bytes(),
		PodID:            common.HexToHash(podId).Bytes(),
	}
}

var checkContractCmd = &cobra.Command{
	Use:   "check",
	Short: "Check whether a payment contract is considered valid",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		validator, err := ethereum.NewPaymentChannelValidator(ethClient, pricingTables, providerAuth)

		_, err = validator.Parse(getPaymentChannelProto(providerAuth, validator.ChainID))
		if err != nil {
			return err
		}

		fmt.Println("Correct Contract")
		return nil
	},
}

var withdrawContractCmd = &cobra.Command{
	Use:   "withdraw <amount>",
	Short: "Withdraw funds from a payment contract",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		providerAuth, _, err := ethereum.GetAccountAndSigner(providerKey, ethClient)
		if err != nil {
			return err
		}

		paymentContract := common.HexToAddress(paymentContractAddress)

		validator, err := ethereum.NewPaymentChannelValidator(
			ethClient,
			map[common.Address]resource.PricingTableMap{paymentContract: make(resource.PricingTableMap)},
			providerAuth,
		)

		channel, err := validator.Parse(getPaymentChannelProto(providerAuth, validator.ChainID))
		if err != nil {
			return err
		}

		_, err = channel.WithdrawUpTo(common.Address{}, big.NewInt(metricsTotal))
		if err != nil {
			return err
		}
		log.Println("Amount transferred successfully")
		return nil
	},
}

func init() {
	contractCmd.AddCommand(checkContractCmd)
	contractCmd.AddCommand(withdrawContractCmd)

	contractCmd.PersistentFlags()

	contractCmd.PersistentFlags().StringVar(&ethereumRpc, "ethereum-rpc", "http://127.0.0.1:8545", "client public address")
	contractCmd.PersistentFlags().StringVar(&providerKey, "ethereum-key", "", "provider account string (private key | http[s]://clef#account | /keystore#account | account (in default keystore))")

	contractCmd.PersistentFlags().StringVar(&paymentContractAddress, "payment-contract", "", "payment contract address")
	contractCmd.PersistentFlags().StringVar(&publisherEthAddress, "publisher", "", "payment contract address")
	contractCmd.PersistentFlags().StringVar(&podId, "pod-id", "00", "pod id")
	withdrawContractCmd.Flags().Int64Var(&metricsTotal, "metric-price", 0, "amount to withdraw up to")
}
