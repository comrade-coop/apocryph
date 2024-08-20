package main

import (
	"fmt"

	"github.com/comrade-coop/apocryph/pkg/publisher"
	"github.com/spf13/cobra"
)

var verifyPodCmd = &cobra.Command{
	Use:     fmt.Sprintf("verify [%s]", publisher.DefaultPodFile),
	Short:   "Verify Pod Images",
	Long:    "Verify the signatures & the certificates of the specified pod images",
	GroupID: "lowlevel",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, pod, _, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}

		err = publisher.VerifyPodImages(pod, publisher.DefaultVerifyOptions(), certificateIdentity, certificateOidcIssuer)
		if err != nil {
			return fmt.Errorf("Failed verifying Pod Images: %v", err)
		}
		return nil
	},
}

func init() {
	verifyPodCmd.Flags().AddFlagSet(verifyImagesFlags)
	podCmd.AddCommand(verifyPodCmd)
}
