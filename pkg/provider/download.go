package provider

import (
	"context"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	iface "github.com/ipfs/boxo/coreiface"
)

func DownloadSecrets(ctx context.Context, ipfs iface.CoreAPI, pod *pb.Pod) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for _, v := range pod.Volumes {
		if v.Type == pb.Volume_VOLUME_SECRET {
			secret := v.GetSecret()
			contents, err := tpipfs.DownloadSecret(ctx, ipfs, secret)
			if err != nil {
				return nil, err
			}
			contents, err = tpipfs.DecryptSecret(secret, contents)
			if err != nil {
				return nil, err
			}
			result[v.Name] = contents
		}
	}
	return result, nil
}

func DownloadImages(ctx context.Context, ipfs iface.CoreAPI, localRegistry string, pod *pb.Pod) (map[string]string, error) {
	result := make(map[string]string)
	for _, c := range pod.Containers {
		ref, err := tpipfs.ReuploadImageFromIpdr(ctx, ipfs, localRegistry, nil, c.Image)
		if err != nil {
			return nil, err
		}
		result[c.Name] = tpipfs.GetRefUrl(ref)
	}
	return result, nil
}
