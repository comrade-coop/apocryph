package main

import (
	"fmt"

	"github.com/comrade-coop/trusted-pods/pkg/abi"
	"github.com/comrade-coop/trusted-pods/pkg/contracts"
	"github.com/comrade-coop/trusted-pods/pkg/crypto"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFileFields = make(map[string]interface{})
var cfgFile string

const (
	keystorePath   string = "/tmp/keystore"
	exportPassword string = "psw"
)

var ProviderAuth *bind.TransactOpts
var Password string = "psw"
var ClientKey string
var PublisherAddress string
var ProviderAddress string
var ClientAddress string
var PaymentContractAddress string
var TokenContractAddress string
var PricePerExecution int64
var MinExecutionPeriod int64
var Units int64
var Instance *abi.Payment
var ethClient *ethclient.Client
var ks *keystore.KeyStore

func setUp() error {

	// create or get a keystore
	ks = crypto.CreateKeystore(keystorePath)
	// setup the payment contract instance
	instance, err := contracts.GetContractInstance(ethClient, common.HexToAddress(PaymentContractAddress))
	Instance = instance
	if err != nil {
		return err
	}

	// get an ethclient
	client, err := contracts.Connect(chainRpc)
	if err != nil {
		return err
	}
	ethClient = client

	// setup a transactor
	_, providerAuth, err := contracts.DeriveAccountConfigs(ProviderKey, Password, exportPassword, client, ks)
	if err != nil {
		return err
	}
	ProviderAuth = providerAuth

	return nil
}

var allowedContractCodeHashes []string = []string{"0x610178dA211FEF7D417bC0e6FeD39F05609AD788, 0x5FbDB2315678afecb367f032d93F642f64180aa3"}

var rootCmd = &cobra.Command{
	Use:   "trustedpods",
	Short: "Trusted Pods is a decentralized compute marketplace for confidential container pods.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use == "help" {
			return nil
		}
		if cfgFile == "" {
			fmt.Fprintln(cmd.ErrOrStderr(), "Warn: config file not set.")
			return nil
		}

		viper.SetConfigFile(cfgFile)

		err := viper.ReadInConfig()
		if err != nil {
			return err
		}

		for k, v := range cfgFileFields {
			viper.UnmarshalKey(k, v)
		}

		return nil
	},
	SilenceUsage: true,
}

func AddConfig(key string, value interface{}, defaultValue interface{}, usage string) {
	viper.SetDefault(key, defaultValue)
	cfgFileFields[key] = value
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	rootCmd.PersistentFlags().StringVar(&chainRpc, "rpc", "http://127.0.0.1:8545", "client public address")
	rootCmd.PersistentFlags().StringVar(&ProviderKey, "key", "", "provider private key")
	rootCmd.PersistentFlags().StringVar(&PaymentContractAddress, "paymentAddr", "", "payment contract address")
	rootCmd.PersistentFlags().StringVar(&ProviderAddress, "provider", "", "provider public address")
	rootCmd.PersistentFlags().Int64Var(&PricePerExecution, "price", 10, "minimum execution price")
	rootCmd.PersistentFlags().Int64Var(&MinExecutionPeriod, "minExecPeriod", 10, "minimum execution period of the pod")

	rootCmd.AddCommand(manifestCmd)
	rootCmd.AddCommand(contractCmd)
	rootCmd.AddCommand(metricsCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(monitorCmd)
}
