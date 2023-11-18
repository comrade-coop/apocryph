package main

import (
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"github.com/comrade-coop/trusted-pods/pkg/ethereum"
	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	"github.com/comrade-coop/trusted-pods/pkg/publisher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// FIXME: Massive code duplication with upload, fund, and sync commands.
var deployPodCmd = &cobra.Command{
	Use:     fmt.Sprintf("deploy [%s] [deployment.yaml]", publisher.DefaultPodFile),
	Aliases: []string{"update"},
	Short:   "Deploy a pod from a local manifest",
	Args:    cobra.MaximumNArgs(2),
	GroupID: "main",
	RunE: func(cmd *cobra.Command, args []string) error {
		podFile, deploymentFile, pod, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		configureDeployment(deployment)

		fundsInt, _ := (&big.Int{}).SetString(funds, 10)
		if fundsInt == nil {
			return fmt.Errorf("Invalid number passed for funds: %s", funds)
		}

		unlockTimeInt := big.NewInt(unlockTime)

		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connecting to IPFS: %w", err)
		}

		if uploadSecrets {
			err = publisher.UploadSecrets(cmd.Context(), ipfs, filepath.Dir(podFile), pod, deployment)
			if err != nil {
				return err
			}
		}

		if uploadImages {
			err = publisher.UploadImages(cmd.Context(), ipfs, pod, deployment)
			if err != nil {
				return err
			}
		}

		err = publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment) // Checkpoint uploads and keys so far
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", err)
		}

		ethClient, err := ethereum.GetClient(ethereumRpc)
		if err != nil {
			return err
		}

		if publisherKey == "" {
			publisherKey = common.BytesToAddress(deployment.Payment.PublisherAddress).String()
		}

		publisherAuth, sign, err := ethereum.GetAccountAndSigner(publisherKey, ethClient)
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

		token := pb.NewToken(string(deployment.Payment.PodID), pb.CreatePod, expirationOffset, publisherAuth.From.Bytes())
		interceptor := &pb.AuthInterceptorClient{Token: token, Sign: sign, ExpirationOffset: time.Duration(expirationOffset) * time.Second}
		err = publisher.SendToProvider(cmd.Context(), tpipfs.NewP2pApi(ipfs, ipfsMultiaddr), pod, deployment, interceptor)
		if err != nil {
			return err
		}

		return publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment)
	},
}

var deletePodCmd = &cobra.Command{
	Use:     fmt.Sprintf("delete [%s|deployment.yaml]", publisher.DefaultPodFile),
	Aliases: []string{"undeploy"},
	Short:   "Delete a pod from a local deployment",
	Args:    cobra.MaximumNArgs(1),
	GroupID: "main",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, deploymentFile, _, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		configureDeployment(deployment)

		ipfs, ipfsMultiaddr, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connecting to IPFS: %w", err)
		}

		err = publisher.SendToProvider(cmd.Context(), tpipfs.NewP2pApi(ipfs, ipfsMultiaddr), nil, deployment)
		if err != nil {
			return err
		}

		return publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment)
	},
}

func init() {
	podCmd.AddCommand(deployPodCmd)
	podCmd.AddCommand(deletePodCmd)

	deployPodCmd.Flags().AddFlagSet(deploymentFlags)
	deployPodCmd.Flags().AddFlagSet(uploadFlags)
	deployPodCmd.Flags().AddFlagSet(fundFlags)
	deployPodCmd.Flags().AddFlagSet(syncFlags)

	deletePodCmd.Flags().AddFlagSet(deploymentFlags)
	deployPodCmd.Flags().AddFlagSet(syncFlags)

}
