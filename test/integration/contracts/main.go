// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/comrade-coop/apocryph/pkg/abi"
	"github.com/comrade-coop/apocryph/pkg/ethereum"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/comrade-coop/apocryph/pkg/resource"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var deadline = big.NewInt(3275538098)
var unlockTime = big.NewInt(5)
var pricePerExecution = big.NewInt(5)
var lockingAmount = big.NewInt(500)
var amountOwed = big.NewInt(5)
var podId = common.HexToHash("00")
var newPrice = big.NewInt(10)

func main() {
	err := mainErr()
	if err != nil {
		fmt.Printf("Error occurred: %v", err)
		os.Exit(1)
	}
}
func mainErr() error {
	if len(os.Args) != 4 {
		return fmt.Errorf("Usage: run-test <PublisherAccountString> <ProviderAccountString> <PaymentContractAddress>")
	}

	ethClient, err := ethereum.GetClient("")
	if err != nil {
		return fmt.Errorf("could not connect to local ethereum node: %w", err)
	}

	publisherAuth, _, err := ethereum.GetAccountAndSigner(os.Args[1], ethClient)
	if err != nil {
		return err
	}

	providerAuth, _, err := ethereum.GetAccountAndSigner(os.Args[2], ethClient)
	if err != nil {
		return err
	}
	paymentAddress := common.HexToAddress(os.Args[3])

	payment, err := abi.NewPayment(paymentAddress, ethClient)
	if err != nil {
		return err
	}

	tokenAddress, err := payment.Token(&bind.CallOpts{})
	if err != nil {
		return err
	}

	token, err := abi.NewMockToken(tokenAddress, ethClient)
	if err != nil {
		return err
	}

	_, err = token.Mint(publisherAuth, big.NewInt(1000))
	if err != nil {
		return fmt.Errorf("could not mint tokens: %w", err)
	}

	_, err = token.Approve(publisherAuth, paymentAddress, lockingAmount)
	if err != nil {
		return err
	}

	fmt.Println("Creating Payment Channel ...")
	_, err = payment.CreateChannel(publisherAuth, providerAuth.From, podId, lockingAmount, unlockTime)
	if err != nil {
		return err
	}

	validator, err := ethereum.NewPaymentChannelValidator(ethClient, map[common.Address]resource.PricingTableMap{paymentAddress: make(resource.PricingTableMap)}, providerAuth)
	if err != nil {
		return err
	}

	channel, err := validator.Parse(&pb.PaymentChannel{
		ChainID:          validator.ChainID.Bytes(),
		ContractAddress:  paymentAddress.Bytes(),
		PublisherAddress: publisherAuth.From.Bytes(),
		ProviderAddress:  providerAuth.From.Bytes(),
		PodID:            podId.Bytes(),
	})
	if err != nil {
		return err
	}

	_, err = channel.WithdrawUpTo(providerAuth.From, amountOwed)
	if err != nil {
		return fmt.Errorf("could not upload metrics and withdraw: %w", err)
	}

	balance, err := token.BalanceOf(&bind.CallOpts{}, providerAuth.From)
	if err != nil {
		return err
	}

	if balance.Cmp(amountOwed) != 0 {
		return fmt.Errorf("Incorrect balance ❌ - expected %v, got %v", amountOwed, balance)
	}
	fmt.Println("Correct balance ✅")

	return nil
}
