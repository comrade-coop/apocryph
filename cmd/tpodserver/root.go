package main

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "trustedpods",
	Short: "Trusted Pods is a decentralized compute marketplace for confidential container pods.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			return errors.New("Config file is required")
		}

		return viper.ReadInConfig()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}
