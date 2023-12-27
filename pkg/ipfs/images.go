// SPDX-License-Identifier: GPL-3.0

package ipfs

import (
	"context"
	"errors"
	"os"
	"strings"

	tpcrypto "github.com/comrade-coop/apocryph/pkg/crypto"
	"github.com/comrade-coop/apocryph/pkg/ipdr"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	imageCopy "github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/docker/daemon"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	encconfig "github.com/containers/ocicrypt/config"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	iface "github.com/ipfs/kubo/core/coreiface"
)

// Convert an [pb.Container_Image] to a [types.ImageReference] plus a digest that can be used for caching later steps.
func GetImageRef(ctx context.Context, sys *types.SystemContext, image *pb.Container_Image) (imageRef types.ImageReference, digest string, err error) {
	if image.Url != "" {
		var imageCloser types.ImageCloser
		imageRef, err = daemon.Transport.ParseReference(image.Url)
		if err == nil {
			imageCloser, err = imageRef.NewImage(ctx, sys)
		}
		if err != nil {
			var err2 error
			imageRef, err2 = docker.Transport.ParseReference("//" + image.Url)
			if err2 == nil {
				imageCloser, err2 = imageRef.NewImage(ctx, sys)
			}
			if err2 != nil {
				err = errors.Join(err, err2)
				return
			}
			err = nil
		}
		defer imageCloser.Close()

		digest = imageCloser.ConfigInfo().Digest.String()
		return
	}
	return
}

// Convert a [types.ImageReference] to a string that can be passed to tools like docker or kubernetes.
func GetRefUrl(imageRef types.ImageReference) string {
	if imageRef.Transport() == docker.Transport {
		return strings.TrimPrefix(imageRef.StringWithinTransport(), "//")
	}
	if imageRef.Transport() == daemon.Transport {
		return imageRef.StringWithinTransport()
	}
	if _, ok := imageRef.Transport().(ipdr.IpdrTransport); ok {
		return imageRef.StringWithinTransport()
	}
	return ""
}

// Encrypt and upload an [types.ImageReference] to IPDR, returning the [pb.Key] and CID (as bytes) of the uploaded image that can then be put into a [pb.Container_Image] or [pb.UploadedImage].
func UploadImageToIpdr(ctx context.Context, ipfs iface.CoreAPI, sys *types.SystemContext, imageRef types.ImageReference) (key *pb.Key, imageCid []byte, err error) {
	ipdrTransport := ipdr.NewIpdrTransport(ipfs)

	copyOptions := &imageCopy.Options{
		DestinationCtx: sys,
		SourceCtx:      sys,
		ReportWriter:   os.Stderr,
	}

	key, err = tpcrypto.NewKey(tpcrypto.KeyTypeOcicrypt)
	if err != nil {
		return
	}

	var cryptoConfig encconfig.CryptoConfig
	cryptoConfig, err = tpcrypto.GetCryptoConfigKey(key)
	if err != nil {
		return
	}

	copyOptions.OciEncryptConfig = cryptoConfig.EncryptConfig
	copyOptions.OciEncryptLayers = &[]int{} // All layers

	policy := &signature.Policy{
		Default: signature.PolicyRequirements{
			signature.NewPRInsecureAcceptAnything(),
		},
	}
	policyContext, _ := signature.NewPolicyContext(policy)
	defer policyContext.Destroy()

	destinationRef := ipdrTransport.NewDestinationReference("")

	_, err = imageCopy.Image(ctx, policyContext, destinationRef, imageRef, copyOptions)
	if err != nil {
		return
	}

	resolved, ok := destinationRef.Path().(path.ImmutablePath)
	if !ok {
		err = errors.New("Destination path not resolved") // Shouldn't get here
		return
	}

	imageCid = resolved.RootCid().Bytes()

	return
}

// Download and decrypt an image from IPDR, reuploading it to a local Docker v2 registry. Takes a [pb.Container_Image] to return the [types.ImageReference] of the uploaded image.
func ReuploadImageFromIpdr(ctx context.Context, ipfs iface.CoreAPI, localRegistryUrl string, sys *types.SystemContext, image *pb.Container_Image) (types.ImageReference, error) {
	if sys == nil {
		sys = &types.SystemContext{
			DockerInsecureSkipTLSVerify: types.OptionalBoolTrue,
		}
	}

	ipdrTransport := ipdr.NewIpdrTransport(ipfs)
	registryTransport := docker.Transport

	if image.Cid != nil && localRegistryUrl != "" {
		copyOptions := &imageCopy.Options{
			DestinationCtx: sys,
			SourceCtx:      sys,
		}

		if image.Key != nil {
			cryptoConfig, err := tpcrypto.GetCryptoConfigKey(image.Key)
			if err != nil {
				return nil, err
			}
			copyOptions.OciDecryptConfig = cryptoConfig.DecryptConfig
		}

		policy := &signature.Policy{
			Default: signature.PolicyRequirements{
				signature.NewPRInsecureAcceptAnything(),
			},
		}
		policyContext, _ := signature.NewPolicyContext(policy)
		defer policyContext.Destroy()

		c, err := cid.Cast(image.Cid)
		if err != nil {
			return nil, err
		}
		sourceRef := ipdrTransport.NewReference(path.FromCid(c), "")

		url := localRegistryUrl + "/" + c.Hash().HexString()
		destinationRef, err := registryTransport.ParseReference("//" + url)
		if err != nil {
			return nil, err
		}

		_, err = imageCopy.Image(ctx, policyContext, destinationRef, sourceRef, copyOptions)

		if err != nil {
			return nil, err
		}

		return destinationRef, nil
	}
	if image.Url != "" {
		ref, _, err := GetImageRef(ctx, sys, image)
		return ref, err
	}
	return nil, errors.New("Failed to read supplied image")
}
