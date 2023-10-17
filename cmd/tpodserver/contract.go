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
	tptypes "github.com/comrade-coop/trusted-pods/pkg/substrate/types"
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
var substrateKey string
var allowedContractCodeHashes []string = []string{"0x610178dA211FEF7D417bC0e6FeD39F05609AD788"}

var Password string = "psw"
var ClientKey string
var ProviderAddress string
var ClientAddress string
var PaymentContractAddress string
var TokenContractAddress string
var PricePerExecution string
var Units int64

const (
	GetSelector    tptypes.ContractSelector = 0x2f865bd9
	ClaimSelector  tptypes.ContractSelector = 0xb388803f
	keystorePath   string                   = "/tmp/keystore"
	exportPassword string                   = "psw"
)

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
		// derive the Account from the private key
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

		_, err = contracts.Withdraw(providerAuth, payment, tokenAddress, clientAddress)
		if err != nil {
			return err
		}
		log.Println("Owed Amount transfered successfully")
		return nil
	},
}
var uploadMetricsCmd = &cobra.Command{
	Use:   "uploadMetric",
	Short: "Upload metric",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		clientAddress := common.HexToAddress(ClientAddress)
		tokenAddress := common.HexToAddress(TokenContractAddress)
		paymentAddress := common.HexToAddress(PaymentContractAddress)
		units := big.NewInt(Units)
		// derive the Account from the private key
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

		_, err = contracts.UploadMetrics(providerAuth, payment, clientAddress, tokenAddress, units)
		if err != nil {
			return err
		}
		log.Println("Uploaded Metrics Successfully")
		return nil
	},
}

func init() {
	contractCmd.AddCommand(checkContractCmd)
	contractCmd.AddCommand(claimCmd)
	contractCmd.AddCommand(uploadMetricsCmd)

	claimCmd.Flags().StringVar(&ClientAddress, "cAddr", "", "client public address")
	claimCmd.Flags().StringVar(&chainRpc, "rpc", "http://127.0.0.1:8545", "client public address")
	claimCmd.Flags().StringVar(&ProviderKey, "key", "", "provider private key")
	claimCmd.Flags().StringVar(&TokenContractAddress, "tokenAddr", "", "token contract address")
	claimCmd.Flags().StringVar(&PaymentContractAddress, "paymentAddr", "", "payment contract address")

	uploadMetricsCmd.Flags().StringVar(&ClientAddress, "cAddr", "", "client public address")
	uploadMetricsCmd.Flags().StringVar(&chainRpc, "rpc", "http://127.0.0.1:8545", "client public address")
	uploadMetricsCmd.Flags().StringVar(&ProviderKey, "key", "", "provider private key")
	uploadMetricsCmd.Flags().StringVar(&TokenContractAddress, "tokenAddr", "", "token contract address")
	uploadMetricsCmd.Flags().StringVar(&PaymentContractAddress, "paymentAddr", "", "payment contract address")
	uploadMetricsCmd.Flags().Int64Var(&Units, "units", 5, "final units of execution")
}
