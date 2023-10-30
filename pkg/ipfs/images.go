package ipfs

import (
	"context"
	"errors"

	"github.com/comrade-coop/trusted-pods/pkg/ipdr"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	imageCopy "github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/docker/daemon"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/go-cid"
)

// TODO: Add way for encrypting(/decrypting) the images

func UploadImagesToIpdr(pod *pb.Pod, ctx context.Context, ipfs iface.CoreAPI, sys *types.SystemContext, keys *[]*pb.Key) error {
	ipdrTransport := ipdr.NewIpdrTransport(ipfs)
	registryTransport := docker.Transport
	daemonTransport := daemon.Transport

	for _, container := range pod.Containers {
		image := container.Image
		if image.Url != "" {
			copyOptions := &imageCopy.Options{}

			policy := &signature.Policy{
				Default: signature.PolicyRequirements{
					signature.NewPRInsecureAcceptAnything(),
				},
			}
			policyContext, _ := signature.NewPolicyContext(policy)
			defer policyContext.Destroy()

			destinationRef := ipdrTransport.NewDestinationReference("")

			sourceRef, err := daemonTransport.ParseReference(image.Url)
			if err != nil {
				return err
			}

			_, err = imageCopy.Image(ctx, policyContext, destinationRef, sourceRef, copyOptions)

			if err != nil {
				sourceRef, err2 := registryTransport.ParseReference("//" + image.Url) // Retry from remote docker registry
				if err2 != nil {
					return errors.Join(err, err2)
				}

				_, err2 = imageCopy.Image(ctx, policyContext, destinationRef, sourceRef, copyOptions)
				if err2 != nil {
					return errors.Join(err, err2)
				}
			}

			resolved, ok := destinationRef.Path().(path.Resolved)
			if !ok {
				return errors.New("Destination path not resolved") // Shouldn't get here
			}

			image.Cid = resolved.Cid().Bytes()
			image.Url = ""
		}
	}
	return nil
}

func ReuploadImagesFromIpdr(pod *pb.Pod, ctx context.Context, ipfs iface.CoreAPI, localRegistryUrl string, sys *types.SystemContext, keys []*pb.Key) error {
	ipdrTransport := ipdr.NewIpdrTransport(ipfs)
	registryTransport := docker.Transport
	for _, container := range pod.Containers {
		image := container.Image
		if len(image.Cid) > 0 {
			copyOptions := &imageCopy.Options{}

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
