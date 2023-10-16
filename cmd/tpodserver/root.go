package main

import (
	"fmt"

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
	rootCmd.AddCommand(manifestCmd)
	rootCmd.AddCommand(contractCmd)
	rootCmd.AddCommand(metricsCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(monitorCmd)
}
