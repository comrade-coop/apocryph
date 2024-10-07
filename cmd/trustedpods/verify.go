package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	tpk8s "github.com/comrade-coop/apocryph/pkg/kubernetes"
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
	Long:    "Verify the signatures & the certificates of the specified image name or Tpod URL",
	Example: "verify ttl.sh/hello-world@sha256:d37ada95d47ad12224c205a938129df7a3e52345828b4fa27b03a98825d1e2e7 --certificate-identity=name@example.com --certificate-oidc-issuer=https://github.com/login/oauth",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parsedURL, err := url.ParseRequestURI(args[0])
		if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
			req, err := http.NewRequest("GET", args[0], nil)
			if err != nil {
				log.Fatalf("Failed to create request: %v", err)
			}
			host := parsedURL.Host
			// Check if the Host is an ip Address by detecting if the port is passed
			if strings.Contains(host, ":") {
				if hostHeader == "" {
					return fmt.Errorf("Must pass the host-header flag when passing an ip endpoint")
				}
				req.Host = hostHeader
			}

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to send request: %v", err)
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("Failed to read response body: %v", err)
			}

			var annotationValues []tpk8s.AnnotationValue
			if err := json.Unmarshal(body, &annotationValues); err != nil {
				return fmt.Errorf("Failed to unmarshal JSON response: %v", err)
			}
			// Verify each image from the response
			verifyOptions := publisher.DefaultVerifyOptions()
			images := []*proto.Image{}
			for _, av := range annotationValues {
				image := &proto.Image{Url: av.URL, VerificationDetails: &proto.VerificationDetails{Signature: av.Signature, Identity: av.Identity, Issuer: av.Issuer}}
				images = append(images, image)
			}
			err = publisher.VerifyImages(images, verifyOptions)
			if err != nil {
				return fmt.Errorf("Failed verifying Images: %v", err)
			}
		} else {
			verifyOptions := publisher.DefaultVerifyOptions()
			if signaturePath != "" {
				verifyOptions.PayloadRef = signaturePath
			}
			image := &proto.Image{Url: args[0], VerificationDetails: &proto.VerificationDetails{Identity: certificateIdentity, Issuer: certificateOidcIssuer}}
			err := publisher.VerifyImages([]*proto.Image{image}, verifyOptions)
			if err != nil {
				return fmt.Errorf("Failed verifying Image: %v", err)
			}
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
