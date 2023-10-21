package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFileFields = make(map[string]interface{})
var cfgFile string

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
}
