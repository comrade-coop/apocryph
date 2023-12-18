// SPDX-License-Identifier: GPL-3.0

package main

import (
	"crypto/sha256"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Operations related to with raw pod manifests",
}

func configureDeployment(deployment *pb.Deployment) error {
	if providerEthAddress != "" && providerPeer != "" {
		deployment.Provider = &pb.ProviderConfig{
			EthereumAddress: common.HexToAddress(providerEthAddress).Bytes(),
			Libp2PAddress:   providerPeer,
		}
	}
	if deployment.Payment == nil {
		deployment.Payment = &pb.PaymentChannelConfig{}
	}
	if paymentContractAddress != "" {
		deployment.Payment.PaymentContractAddress = common.HexToAddress(paymentContractAddress).Bytes()
	}
	if podId != "" {
		deployment.Payment.PodID = common.HexToHash(podId).Bytes()
	}
	if deployment.Payment.PodID == nil {
		podFileNameHash := sha256.Sum256([]byte(deployment.PodManifestFile))
		deployment.Payment.PodID = podFileNameHash[:]
	}
	return nil
}

func init() {
	rootCmd.AddCommand(podCmd)
	rootCmd.AddCommand(registryCmd)

	podCmd.AddGroup(&cobra.Group{
		ID:    "main",
		Title: "Commands",
	})
	podCmd.AddGroup(&cobra.Group{
		ID:    "lowlevel",
		Title: "Low-level commands",
	})

	podCmd.PersistentFlags().AddFlagSet(podFlags)
}
