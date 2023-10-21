package ipfs

import (
	"context"
	"errors"
	"io"
	"os"
	"path"

	tpcrypto "github.com/comrade-coop/trusted-pods/pkg/crypto"
	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	iface "github.com/ipfs/boxo/coreiface"
	ifacepath "github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
)

func UploadSecrets(ctx context.Context, ipfs iface.CoreAPI, basepath string, pod *pb.Pod, keys *[]*pb.Key) error {
	for _, volume := range pod.Volumes {
		if volume.Type == pb.Volume_VOLUME_SECRET {
			err := UploadSecret(ctx, ipfs, basepath, volume.GetSecret(), keys)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UploadSecret(ctx context.Context, ipfs iface.CoreAPI, basepath string, secret *pb.Volume_SecretConfig, keys *[]*pb.Key) error {
	if secret.File != "" {
		secretPath := secret.File
		if !path.IsAbs(secretPath) {
			secretPath = path.Join(basepath, secretPath)
		}
		secretFile, err := os.Open(secretPath)
		if err != nil {
			return err
		}
		defer secretFile.Close()

		secretBytes, err := io.ReadAll(secretFile)
		if err != nil {
			return err
		}

		keyData, err := tpcrypto.CreateRandomKey()
		if err != nil {
			return err
		}
		encryptedSecretBytes, err := tpcrypto.AESEncrypt(secretBytes, keyData)
		if err != nil {
			return err
		}

		encryptedSecretPath, err := ipfs.Unixfs().Add(ctx, files.NewBytesFile(encryptedSecretBytes))
		if err != nil {
			return err
		}

		secret.Cid = encryptedSecretPath.Cid().Bytes()
		secret.KeyIdx = int32(len(*keys))
		secret.File = ""
		*keys = append(*keys, &pb.Key{
			Data: keyData,
		})
	}
	return nil
}

func FetchSecret(ctx context.Context, ipfs iface.CoreAPI, secret *pb.Volume_SecretConfig, keys []*pb.Key) ([]byte, error) {
	var secretBytes []byte
	if secret.Cid != nil {
		secretCid, err := cid.Cast(secret.Cid)
		if err != nil {
			return nil, err
		}
		secretNode, err := ipfs.Unixfs().Get(ctx, ifacepath.IpfsPath(secretCid))
		if err != nil {
			return nil, err
		}
		secretFile, ok := secretNode.(files.File)
		if !ok {
			return nil, errors.New("Supplied secret CID not a file") // TODO: Support encrypted folders
		}
		encryptedSecretBytes, err := io.ReadAll(secretFile)
		if err != nil {
			return nil, err
		}
		keyIdx := secret.KeyIdx
		if keyIdx < 0 {
			secretBytes = encryptedSecretBytes
		} else {
			if int(keyIdx) >= len(keys) {
				return nil, errors.New("Invalid keyIdx")
			}
			key := keys[keyIdx]
			secretBytes, err = tpcrypto.AESDecrypt(encryptedSecretBytes, key.Data)
			if err != nil {
				return nil, err
			}
		}
	}
	return secretBytes, nil
}
