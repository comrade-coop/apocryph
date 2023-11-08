package ipfs

import (
	"context"
	"errors"
	"os"

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

func UploadImagesToIpdr(pod *pb.Pod, ctx context.Context, ipfs iface.CoreAPI, sys *types.SystemContext, keys *[]*pb.Key) error {
	for _, container := range pod.Containers {
		image := container.Image
		if image.Url != "" {
			imageRef, _, err := GetImageRef(ctx, sys, image)
			if err != nil {
				return err
			}

			var key *pb.Key
			key, image.Cid, err = UploadImageToIpdr(ctx, ipfs, sys, imageRef)
			if err != nil {
				return err
			}

			image.KeyIdx = tpcrypto.InsertKey(keys, key)
			image.Url = ""
		}
	}
	return nil
}

func ReuploadImagesFromIpdr(pod *pb.Pod, ctx context.Context, ipfs iface.CoreAPI, localRegistryUrl string, sys *types.SystemContext, keys []*pb.Key) error {
	if sys == nil {
		sys = &types.SystemContext{
			DockerInsecureSkipTLSVerify: types.OptionalBoolTrue,
		}
	}

	ipdrTransport := ipdr.NewIpdrTransport(ipfs)
	registryTransport := docker.Transport
	for _, container := range pod.Containers {
		image := container.Image
		if len(image.Cid) > 0 {
			copyOptions := &imageCopy.Options{
				DestinationCtx: sys,
				SourceCtx:      sys,
			}

			if keys != nil {
				cryptoConfig, err := tpcrypto.GetCryptoConfig(keys, image.KeyIdx)
				if err != nil {
					return err
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
				return err
			}
			sourceRef := ipdrTransport.NewReference(path.IpfsPath(c), "")

			url := localRegistryUrl + "/" + c.Hash().HexString()
			destinationRef, err := registryTransport.ParseReference("//" + url)
			if err != nil {
				return err
			}

			_, err = imageCopy.Image(ctx, policyContext, destinationRef, sourceRef, copyOptions)

			if err != nil {
				return err
			}

			image.Url = url
			image.Cid = nil
		}
	}
	return nil
}
