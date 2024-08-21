package publisher

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/comrade-coop/apocryph/pkg/constants"
	"github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/mitchellh/go-homedir"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/generate"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/sign"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/verify"
	"github.com/spf13/cobra"
)

func DefaultSignOptions() *options.SignOptions {
	cmd := &cobra.Command{}
	o := &options.SignOptions{}
	o.AddFlags(cmd)
	return o
}

func DefaultVerifyOptions() *options.VerifyOptions {
	cmd := &cobra.Command{}
	o := &options.VerifyOptions{}
	o.AddFlags(cmd)
	return o
}

func SignPodImages(pod *proto.Pod, deployment *proto.Deployment, o *options.SignOptions) error {
	var images []string
	for _, container := range pod.Containers {
		images = append(images, container.Image.Url)
	}

	for i, image := range images {
		ro := &options.RootOptions{Timeout: constants.TIMEOUT}

		oidcClientSecret, err := o.OIDC.ClientSecret()
		if err != nil {
			return err
		}
		ko := options.KeyOpts{
			KeyRef:                         o.Key,
			PassFunc:                       generate.GetPass,
			Sk:                             o.SecurityKey.Use,
			Slot:                           o.SecurityKey.Slot,
			FulcioURL:                      o.Fulcio.URL,
			IDToken:                        o.Fulcio.IdentityToken,
			FulcioAuthFlow:                 o.Fulcio.AuthFlow,
			InsecureSkipFulcioVerify:       o.Fulcio.InsecureSkipFulcioVerify,
			RekorURL:                       o.Rekor.URL,
			OIDCIssuer:                     o.OIDC.Issuer,
			OIDCClientID:                   o.OIDC.ClientID,
			OIDCClientSecret:               oidcClientSecret,
			OIDCRedirectURL:                o.OIDC.RedirectURL,
			OIDCDisableProviders:           o.OIDC.DisableAmbientProviders,
			OIDCProvider:                   o.OIDC.Provider,
			SkipConfirmation:               o.SkipConfirmation,
			TSAClientCACert:                o.TSAClientCACert,
			TSAClientCert:                  o.TSAClientCert,
			TSAClientKey:                   o.TSAClientKey,
			TSAServerName:                  o.TSAServerName,
			TSAServerURL:                   o.TSAServerURL,
			IssueCertificateForExistingKey: o.IssueCertificate,
		}

		signaturePath, err := homedir.Expand(constants.OUTPUT_SIGNATURE_PATH)
		if err != nil {
			return err
		}

		err = os.MkdirAll(signaturePath, 0755)
		if err != nil {
			return err
		}

		imageName := strings.ReplaceAll(image, "/", "_")
		signaturePath = signaturePath + "/" + imageName + ".sig"
		_, err = os.Create(signaturePath)
		if err != nil {
			return err
		}

		o.OutputSignature = signaturePath

		if err := sign.SignCmd(ro, ko, *o, []string{image}); err != nil {
			if o.Attachment == "" {
				return fmt.Errorf("signing %v: %w", images, err)
			}
			return fmt.Errorf("signing attachment %s for image %v: %w", o.Attachment, images, err)
		}
		signatureBytes, err := os.ReadFile(signaturePath)
		if err != nil {
			return err
		}
		deployment.Images[i].Signature = string(signatureBytes)
	}
	return nil
}

func VerifyPodImages(pod *proto.Pod, o *options.VerifyOptions, certificateIdentity, certificateOidcIssuer string) error {
	var images []string
	for _, container := range pod.Containers {
		images = append(images, container.Image.Url)
	}

	if o.CommonVerifyOptions.PrivateInfrastructure {
		o.CommonVerifyOptions.IgnoreTlog = true
	}

	annotations, err := o.AnnotationsMap()
	if err != nil {
		return err
	}

	hashAlgorithm, err := o.SignatureDigest.HashAlgorithm()
	if err != nil {
		return err
	}

	o.CertVerify.CertIdentity = certificateIdentity
	o.CertVerify.CertOidcIssuer = certificateOidcIssuer

	v := &verify.VerifyCommand{
		RegistryOptions:              o.Registry,
		CertVerifyOptions:            o.CertVerify,
		CheckClaims:                  o.CheckClaims,
		KeyRef:                       o.Key,
		CertRef:                      o.CertVerify.Cert,
		CertChain:                    o.CertVerify.CertChain,
		CAIntermediates:              o.CertVerify.CAIntermediates,
		CARoots:                      o.CertVerify.CARoots,
		CertGithubWorkflowTrigger:    o.CertVerify.CertGithubWorkflowTrigger,
		CertGithubWorkflowSha:        o.CertVerify.CertGithubWorkflowSha,
		CertGithubWorkflowName:       o.CertVerify.CertGithubWorkflowName,
		CertGithubWorkflowRepository: o.CertVerify.CertGithubWorkflowRepository,
		CertGithubWorkflowRef:        o.CertVerify.CertGithubWorkflowRef,
		IgnoreSCT:                    o.CertVerify.IgnoreSCT,
		SCTRef:                       o.CertVerify.SCT,
		Sk:                           o.SecurityKey.Use,
		Slot:                         o.SecurityKey.Slot,
		Output:                       o.Output,
		RekorURL:                     o.Rekor.URL,
		Attachment:                   o.Attachment,
		Annotations:                  annotations,
		HashAlgorithm:                hashAlgorithm,
		SignatureRef:                 o.SignatureRef,
		PayloadRef:                   o.PayloadRef,
		LocalImage:                   o.LocalImage,
		Offline:                      o.CommonVerifyOptions.Offline,
		TSACertChainPath:             o.CommonVerifyOptions.TSACertChainPath,
		IgnoreTlog:                   o.CommonVerifyOptions.IgnoreTlog,
		MaxWorkers:                   o.CommonVerifyOptions.MaxWorkers,
		ExperimentalOCI11:            o.CommonVerifyOptions.ExperimentalOCI11,
		CertOidcProvider:             o.CertVerify.CertOidcIssuer,
	}

	if o.CommonVerifyOptions.MaxWorkers == 0 {
		return fmt.Errorf("please set the --max-worker flag to a value that is greater than 0")
	}

	if o.Registry.AllowInsecure {
		v.NameOptions = append(v.NameOptions, name.Insecure)
	}

	// if o.CommonVerifyOptions.IgnoreTlog && !o.CommonVerifyOptions.PrivateInfrastructure {
	// 	ui.Warnf(ctx, fmt.Sprintf(ignoreTLogMessage, "signature"))
	// }
	return v.Exec(context.Background(), images)
}
