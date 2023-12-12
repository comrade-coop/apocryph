// SPDX-License-Identifier: GPL-3.0

package ipfs

import (
	"context"
	"errors"
	"os"
	"strings"

	tpcrypto "github.com/comrade-coop/trusted-pods/pkg/crypto"
	"github.com/comrade-coop/trusted-pods/pkg/ipdr"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	imageCopy "github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/docker/daemon"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	encconfig "github.com/containers/ocicrypt/config"
	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/go-cid"
)

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
	copyOptions.OciEncryptLayers = &[]int{}

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

	resolved, ok := destinationRef.Path().(path.Resolved)
	if !ok {
		err = errors.New("Destination path not resolved") // Shouldn't get here
		return
	}

	imageCid = resolved.Cid().Bytes()

	return
}

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
		sourceRef := ipdrTransport.NewReference(path.IpfsPath(c), "")

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
