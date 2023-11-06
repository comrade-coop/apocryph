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

func ReadSecrets(basepath string) SecretTransformation {
	return func(secret *pb.Volume_SecretConfig) error {
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

			secret.Contents = secretBytes
			secret.File = ""
		}
		if secret.ContentsString != "" {
			secret.Contents = []byte(secret.ContentsString)
			secret.ContentsString = ""
		}
		return nil
	}
}

func EncryptSecrets(keys *[]*pb.Key) SecretTransformation {
	return func(secret *pb.Volume_SecretConfig) error {
		var err error
		secret.KeyIdx, err = tpcrypto.CreateKey(keys, tpcrypto.KeyTypeEncrypt)
		if err != nil {
			return err
		}
		secret.Contents, err = tpcrypto.EncryptWith(*keys, secret.KeyIdx, secret.Contents)
		if err != nil {
			return err
		}
		return nil
	}
}

func DecryptSecrets(keys []*pb.Key) SecretTransformation {
	return func(secret *pb.Volume_SecretConfig) error {
		var err error
		secret.Contents, err = tpcrypto.DecryptWith(keys, secret.KeyIdx, secret.Contents)
		if err != nil {
			return err
		}
		return nil
	}
}

func UploadSecrets(ctx context.Context, ipfs iface.CoreAPI) SecretTransformation {
	return func(secret *pb.Volume_SecretConfig) error {
		secretPath, err := ipfs.Unixfs().Add(ctx, files.NewBytesFile(secret.Contents))
		if err != nil {
			return err
		}

		secret.Cid = secretPath.Cid().Bytes()
		secret.Contents = nil
		return nil
	}
}

func DownloadSecrets(ctx context.Context, ipfs iface.CoreAPI) SecretTransformation {
	return func(secret *pb.Volume_SecretConfig) error {
		secretCid, err := cid.Cast(secret.Cid)
		if err != nil {
			return err
		}
		secretNode, err := ipfs.Unixfs().Get(ctx, ifacepath.IpfsPath(secretCid))
		if err != nil {
			return err
		}
		secretFile, ok := secretNode.(files.File)
		if !ok {
			return errors.New("Supplied secret CID not a file") // TODO: Support encrypted folders
		}
		secretBytes, err := io.ReadAll(secretFile)
		if err != nil {
			return err
		}

		secret.Contents = secretBytes
		secret.Cid = nil
		return nil
	}
}