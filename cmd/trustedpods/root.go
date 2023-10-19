package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFileFields = make(map[string]interface{})
var cfgFile string


var PodId string
var Password string = "psw"
var ClientKey string
var ProviderAddress string
var ClientAddress string
var PaymentContractAddress string
var MinAdvanceDuration int64
var TokenContractAddress string
var PricePerExecution int64
var Funds int64
var UnlockTime int64

const (
	keystorePath   string = "/tmp/keystore"
	exportPassword string = "psw"
	rpc            string = "http://127.0.0.1:8545"
)

var allowedContracts []string = []string{"0x5FbDB2315678afecb367f032d93F642f64180aa3"}

var rootCmd = &cobra.Command{
	Use:   "trustedpods",
	Short: "Trusted Pods is a decentralized compute marketplace for confidential container pods.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use == "help" {
			return nil
		}
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath("$XDG_CONFIG_HOME/trustedpods")
			viper.AddConfigPath("$HOME/.config/trustedpods")
			viper.AddConfigPath("$HOME/.trustedpods")
			viper.AddConfigPath("/etc/trustedpods")
		}

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return err
			}
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $XDG_CONFIG_HOME/.trustedpods/config.yaml)")
	rootCmd.AddCommand(podCmd)
	rootCmd.AddCommand(paymentCmd)

	rootCmd.PersistentFlags().StringVar(&ProviderAddress, "provider", "", "provider public address")
	rootCmd.PersistentFlags().StringVar(&ClientAddress, "client", "", "client public address")
	rootCmd.PersistentFlags().StringVar(&PaymentContractAddress, "paymentAddr", "", "payment contract address")
	rootCmd.PersistentFlags().StringVar(&ClientKey, "key", "", "client private key")
	rootCmd.PersistentFlags().StringVar(&Password, "psw", "psw", "password to encrypt the Account")
}
