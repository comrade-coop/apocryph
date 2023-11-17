package main

import (
	"fmt"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pbcon "github.com/comrade-coop/trusted-pods/pkg/proto/protoconnect"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var syncPodCmd = &cobra.Command{
	Use:     fmt.Sprintf("sync [%s|deployment.yaml]", publisher.DefaultPodFile),
	Short:   "Sync a pod from a local deployment",
	Long:    "Sync a pod from a local deployment. Note that this command does not upload any artifacts configured in the pod manifest, as is meant to be used only together with the rest of the low-level commands.",
	GroupID: "lowlevel",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, deploymentFile, pod, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		configureDeployment(deployment)

		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connecting to IPFS: %w", err)
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		if publisherKey == "" {
			publisherKey = common.BytesToAddress(deployment.Payment.PublisherAddress).String()
		}

		_, sign, err := ethereum.GetAccountAndSigner(publisherKey, ethClient)
		if err != nil {
			return fmt.Errorf("Could not get ethereum account: %w", err)
		}

		interceptor := pbcon.NewAuthInterceptorClient(deployment, pbcon.ProvisionPodServiceUpdatePodProcedure, expirationOffset, sign)

		err = publisher.SendToProvider(cmd.Context(), tpipfs.NewP2pApi(ipfs, ipfsMultiaddr), pod, deployment, interceptor, nil, nil)
		if err != nil {
			return err
		}

		return publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment)
	},
}

func init() {
	podCmd.AddCommand(syncPodCmd)

	syncPodCmd.Flags().AddFlagSet(deploymentFlags)
	syncPodCmd.Flags().AddFlagSet(syncFlags)
}
