// SPDX-License-Identifier: GPL-3.0

package provider

import (
	"context"
	"log"

	"github.com/comrade-coop/apocryph/pkg/ipcr"
	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/containerd/containerd"
	iface "github.com/ipfs/kubo/core/coreiface"
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

func DownloadImages(ctx context.Context, client *containerd.Client, ipfsAddress, localRegistry string, pod *pb.Pod) (map[string]string, error) {
	result := make(map[string]string)
	for _, c := range pod.Containers {
		if c.Image.Cid != nil {
			target := c.Image.Url
			if target == "" {
				target = string(c.Image.Cid)
			}
			err := ipcr.PullImage(ctx, client, ipfsAddress, string(c.Image.Cid), target)
			if err != nil {
				return nil, err
			}
			log.Printf("Pulled Image %v Successfully, Decrypting ...\n", c.Image.Url)
			err = ipcr.DecryptImage(ctx, client, "", target, c.Image.Key.Data)
			if err != nil {
				return nil, err
			}
		} else {
			exists, err := ipcr.ImageExists(ctx, client, c.Image.Url)
			if err != nil {
				return nil, err
			}
			if !exists {
				log.Printf("Warning: Image %v does not exist locally.\n", c.Image.Url)
			}
		}
		result[c.Name] = c.Image.Url
	}
	return result, nil
}
