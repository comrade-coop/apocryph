// SPDX-License-Identifier: GPL-3.0

package main

import (
	"fmt"
	"path/filepath"

	"github.com/comrade-coop/apocryph/pkg/ipcr"
	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	"github.com/comrade-coop/apocryph/pkg/publisher"
	"github.com/spf13/cobra"
)

var uploadPodCmd = &cobra.Command{
	Use:     fmt.Sprintf("upload [%s] [deployment.yaml]", publisher.DefaultPodFile),
	Short:   "Upload a pod from a local manifest",
	Long:    "Upload a pod from a local manifest. Note that this command does not send anything to the provider and is meant to be used only together with the rest of the low-level commands.",
	GroupID: "lowlevel",
	Args:    cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		podFile, deploymentFile, pod, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		configureDeployment(deployment)

		ipfs, _, err := tpipfs.GetIpfsClient(ipfsApi)
		if err != nil {
			return fmt.Errorf("Failed connecting to IPFS: %w", err)
		}

		if signImages {
			err := checkCertificateFlags()
			if err != nil {
				return err
			}
			signOptions := publisher.DefaultSignOptions()
			if !uploadSignatures {
				signOptions.Upload = false
			}
			err = publisher.SignPodImages(pod, deployment, signOptions, certificateIdentity, certificateOidcIssuer)
			if err != nil {
				return fmt.Errorf("failed Signing images: %v", err)
			}
		}

		ctrdClient, err := ipcr.GetContainerdClient("k8s.io")
		if err != nil {
			return err
		}

		if uploadSecrets {
			err = publisher.UploadSecrets(cmd.Context(), ipfs, filepath.Dir(podFile), pod, deployment)
			if err != nil {
				return err
			}
		}

		if uploadImages {
			err = publisher.UploadImages(cmd.Context(), ctrdClient, ipfsApi, pod, deployment)
			if err != nil {
				return err
			}
		}

		return publisher.SaveDeployment(deploymentFile, deploymentFormat, deployment)
	},
}

func init() {
	podCmd.AddCommand(uploadPodCmd)

	uploadPodCmd.Flags().AddFlagSet(deploymentFlags)
	uploadPodCmd.Flags().AddFlagSet(uploadFlags)
	uploadPodCmd.Flags().AddFlagSet(imageCertificateFlags)
}
