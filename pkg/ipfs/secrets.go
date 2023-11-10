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

type SecretTransformation func(secret *pb.Volume_SecretConfig) error

func TransformSecrets(pod *pb.Pod, transformations ...SecretTransformation) error {
	for _, volume := range pod.Volumes {
		if volume.Type == pb.Volume_VOLUME_SECRET {
			for _, transformation := range transformations {
				err := transformation(volume.GetSecret())
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}

func ReadSecret(basepath string, secret *pb.Volume_SecretConfig) ([]byte, error) {
	if secret.File != "" {
		secretPath := secret.File
		if !path.IsAbs(secretPath) {
			secretPath = path.Join(basepath, secretPath)
		}
		secretFile, err := os.Open(secretPath)
		if err != nil {
			return nil, err
		}
		defer secretFile.Close()

		secretBytes, err := io.ReadAll(secretFile)
		if err != nil {
			return nil, err
		}

		return secretBytes, nil
	}
	if secret.ContentsString != "" {
		return []byte(secret.ContentsString), nil
	}
	return nil, nil
}

func EncryptSecret(data []byte) (key *pb.Key, contents []byte, err error) {
	key, err = tpcrypto.NewKey(tpcrypto.KeyTypeEncrypt)
	if err != nil {
		return
	}
	contents, err = tpcrypto.EncryptWithKey(key, data)
	return
}

func UploadSecret(ctx context.Context, ipfs iface.CoreAPI, contents []byte) (cid []byte, err error) {
	secretPath, err := ipfs.Unixfs().Add(ctx, files.NewBytesFile(contents))
	if err != nil {
		return nil, err
	}
	return secretPath.Cid().Bytes(), nil
}

func DownloadSecret(ctx context.Context, ipfs iface.CoreAPI, secret *pb.Volume_SecretConfig) ([]byte, error) {
	if secret.Contents != nil {
		return secret.Contents, nil
	}
	secretCid, err := cid.Cast(secret.Cid)
	if err != nil {
		return nil, err
	}
	secretNode, err := ipfs.Unixfs().Get(ctx, ifacepath.IpfsPath(secretCid))
	if err != nil {
		return nil, err
	}
	defer secretNode.Close()
	secretFile, ok := secretNode.(files.File)
	if !ok {
		return nil, errors.New("Supplied secret CID not a file") // TODO: Support encrypted folders
	}
	return io.ReadAll(secretFile)
}

func DecryptSecret(secret *pb.Volume_SecretConfig, contents []byte) ([]byte, error) {
	if secret.Key == nil {
		return contents, nil
	}
	return tpcrypto.DecryptWithKey(secret.Key, contents)
}
