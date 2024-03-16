// SPDX-License-Identifier: GPL-3.0

package publisher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/comrade-coop/apocryph/pkg/ipcr"
	pb "github.com/comrade-coop/apocryph/pkg/proto"
	"github.com/containerd/containerd"
	iface "github.com/ipfs/kubo/core/coreiface"
	"golang.org/x/exp/slices"

	tpipfs "github.com/comrade-coop/apocryph/pkg/ipfs"
)

func UploadSecrets(ctx context.Context, ipfs iface.CoreAPI, basepath string, pod *pb.Pod, deployment *pb.Deployment) error {
	fmt.Fprintf(os.Stderr, "Encrypting and uploading secrets to IPFS...\n")

	oldUploadedSecrets := map[string]*pb.UploadedSecret{}
	for _, uploadedSecret := range deployment.Secrets {
		oldUploadedSecrets[uploadedSecret.VolumeName] = uploadedSecret
	}
	deployment.Secrets = []*pb.UploadedSecret{}

	for _, volume := range pod.Volumes {
		if volume.Type == pb.Volume_VOLUME_SECRET {
			secretBytes, err := tpipfs.ReadSecret(basepath, volume.GetSecret())
			if err != nil {
				return err
			}

			secretSha256 := sha256.Sum256(secretBytes)
			uploadedSecret, ok := oldUploadedSecrets[volume.Name]
			if !ok || !slices.Equal(uploadedSecret.Sha256Sum, secretSha256[:]) {
				if ok {
					err = tpipfs.RemoveSecret(ctx, ipfs, uploadedSecret.Cid)
					if err != nil {
						fmt.Printf("Failed unpinning old secret from IPFS: %v", err)
					}
				}
				// Have to reupload...
				key, secretContents, err := tpipfs.EncryptSecret(secretBytes)
				if err != nil {
					return fmt.Errorf("Failed uploading secrets to IPFS: %w", err)
				}
				secretCid, err := tpipfs.UploadSecret(ctx, ipfs, secretContents)
				if err != nil {
					return fmt.Errorf("Failed uploading secrets to IPFS: %w", err)
				}

				uploadedSecret = &pb.UploadedSecret{
					VolumeName: volume.Name,
					Sha256Sum:  secretSha256[:],
					Cid:        secretCid,
					Key:        key,
				}
				fmt.Fprintf(os.Stderr, "Encrypted and uploaded secret `%s` to IPFS\n", volume.Name)
			}

			deployment.Secrets = append(deployment.Secrets, uploadedSecret)
		}
	}
	return nil
}

func UploadImages(ctx context.Context, client *containerd.Client, IPFSAddress string, pod *pb.Pod, deployment *pb.Deployment) error {
	fmt.Fprintf(os.Stderr, "Encrypting and uploading images to IPCR...\n")

	oldUploadedImages := map[string]*pb.UploadedImage{}
	for _, uploadedImage := range deployment.Images {
		oldUploadedImages[uploadedImage.SourceUrl] = uploadedImage
	}
	deployment.Images = []*pb.UploadedImage{}

	for _, container := range pod.Containers {
		image := container.Image

		err := ipcr.EnsureImage(ctx, client, image.Url)
		if err != nil {
			return err
		}

		_, prvKey, err := ipcr.EncryptImage(ctx, client, image.Url, "")
		if err != nil {
			return err
		}

		cid, err := ipcr.PushImage(ctx, client, IPFSAddress, image.Url)
		if err != nil {
			return err
		}

		uploadedImage := &pb.UploadedImage{SourceUrl: image.Url, Cid: []byte(cid), Key: &pb.Key{Data: prvKey}}
		deployment.Images = append(deployment.Images, uploadedImage)
	}
	return nil
}

func LinkUploadsFromDeployment(pod *pb.Pod, deployment *pb.Deployment) *pb.Pod {
	if pod == nil {
		return nil
	}

	uploadedSecrets := map[string]*pb.UploadedSecret{}
	for _, uploadedSecret := range deployment.Secrets {
		uploadedSecrets[uploadedSecret.VolumeName] = uploadedSecret
	}

	for _, volume := range pod.Volumes {
		if volume.Type == pb.Volume_VOLUME_SECRET {
			if uploadedSecret, ok := uploadedSecrets[volume.Name]; ok {
				volume.Configuration = &pb.Volume_Secret{
					Secret: &pb.Volume_SecretConfig{
						Cid: uploadedSecret.Cid,
						Key: uploadedSecret.Key,
					},
				}
			}
		}
	}

	uploadedImages := map[string]*pb.UploadedImage{}
	for _, uploadedImage := range deployment.Images {
		uploadedImages[uploadedImage.SourceUrl] = uploadedImage
	}
	deployment.Images = []*pb.UploadedImage{}

	for _, container := range pod.Containers {
		if uploadedImage, ok := uploadedImages[container.Image.Url]; ok {
			container.Image = &pb.Container_Image{
				Cid: uploadedImage.Cid,
				Url: uploadedImage.SourceUrl,
				Key: uploadedImage.Key,
			}
		}
	}

	return pod
}
