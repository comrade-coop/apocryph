package main

import (
	"log"
	"math/big"
	"os"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/comrade-coop/trusted-pods/pkg/contracts"
	"github.com/comrade-coop/trusted-pods/pkg/crypto"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var deadline = big.NewInt(3275538098)
var mindAdvanceDuration = big.NewInt(5)
var pricePerExecution = big.NewInt(5)
var lockingAmount = big.NewInt(500)
var unitsOfExecution = big.NewInt(5)
var newPrice = big.NewInt(10)

func main() {
	client, provider, payment, token, tokenAddress, paymentAddress, clientAuth, providerAuth, c := setUp()
	token.Approve(clientAuth, paymentAddress, lockingAmount)

	log.Println("Creating Payment Channel ...")
	_, err := contracts.CreatePaymentChannel(clientAuth, payment, provider.Address, tokenAddress, lockingAmount, deadline, mindAdvanceDuration, pricePerExecution)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}

	log.Println("Uploading Metrics ...")
	_, err = contracts.UploadMetrics(providerAuth, payment, client.Address, tokenAddress, unitsOfExecution)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}
	log.Println("withdrawing owed Amount")
	_, err = contracts.Withdraw(providerAuth, payment, tokenAddress, client.Address)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}

	balance, err := contracts.Balance(clientAuth, c, token, provider.Address)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}
	log.Printf("Provider Balance: %v", balance)

	_, err = contracts.UpdatePrice(providerAuth, payment, client.Address, tokenAddress, newPrice)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}
	log.Println("Suggested new Price:", newPrice)

	_, err = contracts.AcceptPrice(providerAuth, payment, provider.Address, tokenAddress)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}
	log.Println("New Price Accepted")

	log.Println("Uploading Metrics ...")
	_, err = contracts.UploadMetrics(providerAuth, payment, client.Address, tokenAddress, unitsOfExecution)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}

	_, err = contracts.Withdraw(providerAuth, payment, tokenAddress, client.Address)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}

	balance, err = contracts.Balance(clientAuth, c, token, provider.Address)
	if err != nil {
		log.Printf("Error occured: %v", err)
	}

	log.Printf("Provider Balance: %v", balance)
	if balance.Int64() == 75 {
		log.Println("Correct expected balance ✅")
	} else {
		log.Println("Incorrect balance ❌")
	}

}

func setUp() (accounts.Account, accounts.Account, *abi.Payment, *abi.MockToken, common.Address, common.Address, *bind.TransactOpts, *bind.TransactOpts, *ethclient.Client) {
	log.Println("Creating Provider & Client Accounts ...")
	args := os.Args
	if len(args) != 4 {
		log.Fatalf("Usage: run-test <KesytorePath> <ClientPrivateKey> <ProviderPrivateKey>")
	}
	ks := crypto.CreateKeystore(args[1])
	key, err := crypto.DerivePrvKey(args[2][2:])
	if err != nil {
		log.Fatalf("could not derive private key: %v", err)
	}
	clientAcc, err := crypto.AccountFromECDSA(key, "psw", ks)
	if err != nil {
		log.Printf("error creating account: %v \n", err)
	}
	keyjson, err := crypto.ExportAccount(clientAcc, "psw", "psw", ks)
	if err != nil {
		log.Printf("error exporting account: %v", err)
	}

	providerKey, err := crypto.DerivePrvKey(args[3][2:])
	if err != nil {
		log.Fatalf("could not derive private key: %v", err)
	}
	providerAcc, err := crypto.AccountFromECDSA(providerKey, "psw", ks)
	if err != nil {
		log.Printf("error creating account: %v \n", err)
	}

	providerKeyjson, err := crypto.ExportAccount(providerAcc, "psw", "psw", ks)
	if err != nil {
		log.Printf("error exporting account: %v", err)
	}

	client, err := contracts.ConnectToLocalNode()
	if err != nil {
		log.Fatalf("could not connect to local ethereum node")
	}
	clientAuth, err := contracts.CreateTransactor(string(keyjson), "psw", client)
	providerAuth, err := contracts.CreateTransactor(string(providerKeyjson), "psw", client)
	log.Println("Accounts Retreived")
	paymentAddress, payment, _ := contracts.DeployPaymentContract(clientAuth, client)
	tokenAddress, token, _ := contracts.DeployTokenContract(clientAuth, client)
	_, err = contracts.ClaimTokens(clientAuth, client, token, big.NewInt(1000))
	if err != nil {
		log.Fatalf("could not claim tokens: %v", err)
	}
	_, err = contracts.Balance(clientAuth, client, token, clientAcc.Address)
	if err != nil {
		log.Fatalf("could not get balance: %v", err)
	}
	return clientAcc, providerAcc, payment, token, *tokenAddress, *paymentAddress, clientAuth, providerAuth, client
}
