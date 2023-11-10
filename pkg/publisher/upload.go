package publisher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	iface "github.com/ipfs/boxo/coreiface"
	"golang.org/x/exp/slices"

	tpipfs "github.com/comrade-coop/trusted-pods/pkg/ipfs"
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

func UploadImages(ctx context.Context, ipfs iface.CoreAPI, pod *pb.Pod, deployment *pb.Deployment) error {
	fmt.Fprintf(os.Stderr, "Encrypting and uploading images to IPDR...\n")

	oldUploadedImages := map[string]*pb.UploadedImage{}
	for _, uploadedImage := range deployment.Images {
		oldUploadedImages[uploadedImage.SourceUrl] = uploadedImage
	}
	deployment.Images = []*pb.UploadedImage{}

	for _, container := range pod.Containers {
		image := container.Image
		imageRef, imageDigest, err := tpipfs.GetImageRef(ctx, nil, image)
		if err != nil {
			return fmt.Errorf("Failed parsing images reference: %w", err)
		}
		uploadedImage, ok := oldUploadedImages[image.Url]
		if !ok || imageDigest != uploadedImage.Digest {
			fmt.Fprintf(os.Stderr, "Uploading image `%s`...\n", image.Url)
			imageKey, imageCid, err := tpipfs.UploadImageToIpdr(ctx, ipfs, nil, imageRef)
			if err != nil {
				return fmt.Errorf("Failed encrypting and uploading image to IPDR: %w", err)
			}

			uploadedImage = &pb.UploadedImage{
				SourceUrl: image.Url,
				Digest:    imageDigest,
				Cid:       imageCid,
				Key:       imageKey,
			}
			fmt.Fprintf(os.Stderr, "Uploaded image `%s`\n", image.Url)
		}

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
				Key: uploadedImage.Key,
			}
		}
	}

	return pod
}
