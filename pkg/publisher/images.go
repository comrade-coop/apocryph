package publisher

import (
	"fmt"
	"time"

	"github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/generate"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/sign"
	"github.com/spf13/cobra"
)

func DefaultSignOptions() *options.SignOptions {
	cmd := &cobra.Command{}
	o := &options.SignOptions{}
	o.AddFlags(cmd)
	return o
}

func SignPodImages(pod *proto.Pod, o *options.SignOptions) error {
	var images []string
	for _, container := range pod.Containers {
		images = append(images, container.Image.Url)
	}

	ro := &options.RootOptions{Timeout: 3 * time.Minute}

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
	if err := sign.SignCmd(ro, ko, *o, images); err != nil {
		if o.Attachment == "" {
			return fmt.Errorf("signing %v: %w", images, err)
		}
		return fmt.Errorf("signing attachment %s for image %v: %w", o.Attachment, images, err)
	}
	return nil
}
