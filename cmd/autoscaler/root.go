// SPDX-License-Identifier: GPL-3.0

package main

import (
	"github.com/spf13/cobra"
)

const (
	keystorePath string = "/tmp/keystore"
)

var rootCmd = &cobra.Command{
	Use:   "Autonomous Autoscaler",
	Short: "Autoscaler is autonomous application offering service discovery and fault-tolerance of applications accross Apocryph network",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(startCmd)
}
