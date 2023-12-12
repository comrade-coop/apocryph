// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"math/big"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var fundPodCmd = &cobra.Command{
	Use:     fmt.Sprintf("fund [%s|deployment.yaml]", publisher.DefaultPodFile),
	Short:   "Fund the pod's execution",
	Args:    cobra.MaximumNArgs(1),
	GroupID: "main",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, deploymentFile, _, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		configureDeployment(deployment)

		fundsInt, _ := (&big.Int{}).SetString(funds, 10)
		if fundsInt == nil {
			return fmt.Errorf("Invalid number passed for funds: %s", funds)
		}

		unlockTimeInt := big.NewInt(unlockTime)

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		if publisherKey == "" {
			publisherKey = common.BytesToAddress(deployment.Payment.PublisherAddress).String()
		}

		publisherAuth, _, err := ethereum.GetAccountAndSigner(publisherKey, ethClient)
		if err != nil {
			return fmt.Errorf("Could not get ethereum account: %w", err)
		}

		// FIXME: move to configureDeployment?
		chainId, err := ethClient.ChainID(cmd.Context())
		if err != nil {
			return err
		}
		deployment.Payment.ChainID = chainId.Bytes()
		deployment.Payment.PublisherAddress = publisherAuth.From.Bytes()

		err = publisher.FundPaymentChannel(ethClient, publisherAuth, deployment, fundsInt, unlockTimeInt, debugMintFunds)
		if err != nil {
			return err
		}

		return publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment)
	},
}

func init() {
	podCmd.AddCommand(fundPodCmd)

	fundPodCmd.Flags().AddFlagSet(fundFlags)
	fundPodCmd.Flags().AddFlagSet(deploymentFlags)

}
