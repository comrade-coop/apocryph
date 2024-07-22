package ipcr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tpcrypto "github.com/comrade-coop/apocryph/pkg/crypto"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/images/converter"
	"github.com/containerd/log"
	"github.com/containerd/nerdctl/pkg/api/types"
	img "github.com/containerd/nerdctl/pkg/cmd/image"
	"github.com/containerd/nerdctl/pkg/imgutil"
	"github.com/containerd/nerdctl/pkg/ipfs"
	"github.com/containerd/nerdctl/pkg/referenceutil"
	encutils "github.com/containers/ocicrypt/utils"
	"github.com/spf13/cobra"
)

const RSA_KEY_SIZE = 2048

var cryptOptions = types.ImageCryptOptions{
	Stdout:        os.Stdout,
	Platforms:     nil,
	AllPlatforms:  false,
	GpgHomeDir:    "",
	GpgVersion:    "",
	Keys:          nil,
	DecRecipients: nil,
	Recipients:    []string{"jwe:pubkey.pem)"},
}

// EncryptImage returns the private key for decryption
// password could be ommitted
func EncryptImage(ctx context.Context, client *containerd.Client, image, password string) ([]byte, []byte, error) {
	pubKey, prvKey, err := encutils.CreateRSATestKey(RSA_KEY_SIZE, []byte(password), true)
	if err != nil {
		return nil, nil, err
	}
	err = tpcrypto.Crypt(ctx, client, image, image, true, cryptOptions, [][]byte{pubKey}, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	return pubKey, prvKey, nil
}
func DecryptImage(ctx context.Context, client *containerd.Client, password, image string, prvKey []byte) error {
	err := tpcrypto.Crypt(ctx, client, image, image, false, cryptOptions, nil, [][]byte{prvKey}, [][]byte{[]byte(password)})
	if err != nil {
		return err
	}
	return nil
}
func PushImage(ctx context.Context, client *containerd.Client, ipfsAddr, image string) (string, error) {

	cmd := cobra.Command{}
	options := types.ImagePushOptions{
		GOptions:     types.GlobalCommandOptions{},
		AllPlatforms: false,
		Quiet:        false,
		Stdout:       cmd.OutOrStdout(),
		IpfsAddress:  ipfsAddr,
	}
	c, err := IpfsPush(ctx, client, "ipfs://"+image, options)
	if err != nil {
		return "", err
	}
	return c, nil
}

// copied from containerd/nerdctl/ipfs + returning the cid
func IpfsPush(ctx context.Context, client *containerd.Client, rawRef string, options types.ImagePushOptions) (string, error) {
	if scheme, ref, err := referenceutil.ParseIPFSRefWithScheme(rawRef); err == nil {
		if scheme != "ipfs" {
			return "", fmt.Errorf("ipfs scheme is only supported but got %q", scheme)
		}
		log.G(ctx).Infof("pushing image %q to IPFS", ref)

		var ipfsPath string
		if options.IpfsAddress != "" {
			dir, err := os.MkdirTemp("", "apidirtmp")
			if err != nil {
				return "", err
			}
			defer os.RemoveAll(dir)
			if err := os.WriteFile(filepath.Join(dir, "api"), []byte(options.IpfsAddress), 0600); err != nil {
				return "", err
			}
			ipfsPath = dir
		}
		// removed the condition here
		var layerConvert converter.ConvertFunc

		c, err := ipfs.Push(ctx, client, ref, layerConvert, options.AllPlatforms, options.Platforms, options.IpfsEnsureImage, ipfsPath)
		if err != nil {
			log.G(ctx).WithError(err).Warnf("ipfs push failed")
			return "", err
		}
		fmt.Fprintln(options.Stdout, c)
		return c, nil
	}
	return "", fmt.Errorf("Could not parse ipfs image name")
}

// EnsureImage makes sure the image is available in containerd store
func EnsureImage(ctx context.Context, client *containerd.Client, image string) error {
	options := types.ImagePullOptions{
		GOptions:      types.GlobalCommandOptions{},
		VerifyOptions: types.ImageVerifyOptions{Provider: "none"},
		AllPlatforms:  false,
		Platform:      nil,
		Unpack:        "false",
		Quiet:         false,
		RFlags: imgutil.RemoteSnapshotterFlags{
			SociIndexDigest: "",
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err := img.Pull(ctx, client, image, options)
	if err != nil {
		return err
	}
	return nil
}

func ImageExists(ctx context.Context, client *containerd.Client, image string) (bool, error) {
	images, err := img.List(ctx, client, []string{"reference=" + image}, nil)
	if err != nil {
		return false, err
	}
	for _, im := range images {
		if im.Name == image {
			return true, nil
		}
	}
	return false, nil
}

func PullImage(ctx context.Context, client *containerd.Client, ipfsAddr, image, target string) error {
	cmd := cobra.Command{}
	options := types.ImagePullOptions{
		GOptions:      types.GlobalCommandOptions{},
		VerifyOptions: types.ImageVerifyOptions{Provider: "none"},
		AllPlatforms:  false,
		Platform:      nil,
		Unpack:        "false",
		Quiet:         false,
		IPFSAddress:   ipfsAddr,
		RFlags: imgutil.RemoteSnapshotterFlags{
			SociIndexDigest: "",
		},
		Stdout: cmd.OutOrStdout(),
		Stderr: cmd.OutOrStderr(),
	}
	err := img.Pull(ctx, client, "ipfs://"+image, options)
	if err != nil {
		return err
	}

	err = img.Tag(ctx, client, types.ImageTagOptions{Source: image, Target: target})
	if err != nil {
		return err
	}

	return nil
}

func GetContainerdClient(namespace string) (*containerd.Client, error) {
	client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace(namespace))
	if err != nil {
		return nil, err
	}
	return client, nil
}
