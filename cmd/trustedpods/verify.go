package main

import (
	"fmt"

	"github.com/comrade-coop/apocryph/pkg/proto"
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
		_, _, pod, deployment, err := publisher.ReadPodAndDeployment(args, manifestFormat, deploymentFormat)
		if err != nil {
			return err
		}
		// in case user deployed the pod himself and wants to verify it and the
		// pod does not already specify a certificate identity & owner
		pod = publisher.LinkUploadsFromDeployment(pod, deployment)
		err = publisher.VerifyPodImages(pod, publisher.DefaultVerifyOptions())
		if err != nil {
			return fmt.Errorf("Failed verifying Pod Images: %v", err)
		}
		return nil
	},
}

var verifyImageCmd = &cobra.Command{
	Use:     fmt.Sprintf("verify image"),
	Short:   "Verify image signature",
	Long:    "Verify the signatures & the certificates of the specified image name",
	Example: "verify ttl.sh/hello-world@sha256:d37ada95d47ad12224c205a938129df7a3e52345828b4fa27b03a98825d1e2e7 --certificate-identity=name@example.com --certificate-oidc-issuer=https://github.com/login/oauth",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		verifyOptions := publisher.DefaultVerifyOptions()
		if signaturePath != "" {
			verifyOptions.PayloadRef = signaturePath
		}
		image := &proto.Image{Url: args[0], VerificationDetails: &proto.VerificationDetails{Identity: certificateIdentity, Issuer: certificateOidcIssuer}}
		err := publisher.VerifyImages([]*proto.Image{image}, verifyOptions)
		if err != nil {
			return fmt.Errorf("Failed verifying Image: %v", err)
		}

		return nil
	},
}

func init() {
	verifyPodCmd.Flags().AddFlagSet(verifyFlags)
	verifyImageCmd.Flags().AddFlagSet(verifyFlags)
	podCmd.AddCommand(verifyPodCmd)
	rootCmd.AddCommand(verifyImageCmd)
}
