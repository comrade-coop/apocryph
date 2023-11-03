package publisher

import (
	"context"
	"fmt"
	"os"
	"path"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
)

func UploadManifest(ctx context.Context, manifestFile, manifestFormat, ipfsApi string, noIpdr bool) (*pb.ProvisionPodRequest, *pb.Pod, error) {
	ipfs, _, err := tpipfs.GetIpfsClient(ipfsApi)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed connectig to IPFS: %w", err)
	}

	pod := &pb.Pod{}
	err = pb.UnmarshalFile(manifestFile, manifestFormat, pod)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed reading the manifest file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Parsed pod manifest\n")

	keys := []*pb.Key{}

	err = tpipfs.TransformSecrets(pod,
		tpipfs.ReadSecrets(path.Dir(manifestFile)),
		tpipfs.EncryptSecrets(&keys),
		tpipfs.UploadSecrets(ctx, ipfs),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed encrypting and uploading secrets to IPFS: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Encrypted and uploaded pod secrets to IPFS\n")

	if !noIpdr {
		err = tpipfs.UploadImagesToIpdr(pod, ctx, ipfs, nil, &keys)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed encrypting and uploading images to IPDR: %w", err)
		}
	}

	fmt.Fprintf(os.Stderr, "Encrypted and uploaded pod images to IPDR\n")

	podCid, err := tpipfs.AddProtobufFile(ipfs, pod)
	if err != nil {
		return nil, nil, err
	}

	fmt.Fprintf(os.Stderr, "Uploaded pod manifest to IPFS\n")

	return &pb.ProvisionPodRequest{
		PodManifestCid: podCid.Bytes(),
		Keys:           keys,
	}, pod, nil
}
