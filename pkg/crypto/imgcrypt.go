package crypto

import (
	"context"
	"errors"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images/converter"
	"github.com/containerd/imgcrypt/images/encryption"
	"github.com/containerd/imgcrypt/images/encryption/parsehelpers"
	"github.com/containerd/nerdctl/pkg/api/types"
	"github.com/containerd/nerdctl/pkg/platformutil"
	"github.com/containerd/nerdctl/pkg/referenceutil"
	encconfig "github.com/containers/ocicrypt/config"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// from containerd/nerdctl/pkg/cmd/image/crypt.go modified to take keys directly
// instead of reading them from the file system
// args used for ecryption/decryption
func Crypt(ctx context.Context, client *containerd.Client, srcRawRef, targetRawRef string, encrypt bool, options types.ImageCryptOptions, pubKeys, prvKeys [][]byte, privKeysPasswords [][]byte) error {
	var convertOpts = []converter.Opt{}
	if srcRawRef == "" || targetRawRef == "" {
		return errors.New("src and target image need to be specified")
	}

	srcNamed, err := referenceutil.ParseAny(srcRawRef)
	if err != nil {
		return err
	}
	srcRef := srcNamed.String()

	targetNamed, err := referenceutil.ParseDockerRef(targetRawRef)
	if err != nil {
		return err
	}
	targetRef := targetNamed.String()

	platMC, err := platformutil.NewMatchComparer(options.AllPlatforms, options.Platforms)
	if err != nil {
		return err
	}
	convertOpts = append(convertOpts, converter.WithPlatform(platMC))

	imgcryptFlags, err := parseImgcryptFlags(options, encrypt)
	if err != nil {
		return err
	}

	srcImg, err := client.ImageService().Get(ctx, srcRef)
	if err != nil {
		return err
	}
	layerDescs, err := platformutil.LayerDescs(ctx, client.ContentStore(), srcImg.Target, platMC)
	if err != nil {
		return err
	}
	layerFilter := func(desc ocispec.Descriptor) bool {
		return true
	}
	var convertFunc converter.ConvertFunc
	if encrypt {
		cc, err := CreateCryptoConfig(imgcryptFlags, layerDescs, pubKeys)
		if err != nil {
			return err
		}
		convertFunc = encryption.GetImageEncryptConverter(&cc, layerFilter)
	} else {
		cc, err := CreateDecryptCryptoConfig(imgcryptFlags, layerDescs, privKeysPasswords, prvKeys)
		if err != nil {
			return err
		}
		convertFunc = encryption.GetImageDecryptConverter(&cc, layerFilter)
	}
	// we have to compose the DefaultIndexConvertFunc here to match platforms.
	convertFunc = composeConvertFunc(converter.DefaultIndexConvertFunc(nil, false, platMC), convertFunc)
	convertOpts = append(convertOpts, converter.WithIndexConvertFunc(convertFunc))

	// converter.Convert() gains the lease by itself
	newImg, err := converter.Convert(ctx, client, targetRef, srcRef, convertOpts...)
	if err != nil {
		return err
	}
	fmt.Fprintln(options.Stdout, newImg.Target.Digest.String())
	return nil
}

func composeConvertFunc(a, b converter.ConvertFunc) converter.ConvertFunc {
	return func(ctx context.Context, cs content.Store, desc ocispec.Descriptor) (*ocispec.Descriptor, error) {
		newDesc, err := a(ctx, cs, desc)
		if err != nil {
			return newDesc, err
		}
		if newDesc == nil {
			return b(ctx, cs, desc)
		}
		return b(ctx, cs, *newDesc)
	}
}

// parseImgcryptFlags corresponds to https://github.com/containerd/imgcrypt/blob/v1.1.2/cmd/ctr/commands/images/crypt_utils.go#L244-L252
func parseImgcryptFlags(options types.ImageCryptOptions, encrypt bool) (parsehelpers.EncArgs, error) {
	var a parsehelpers.EncArgs

	a.GPGHomedir = options.GpgHomeDir
	a.GPGVersion = options.GpgVersion
	a.Key = options.Keys
	if encrypt {
		a.Recipient = options.Recipients
		if len(a.Recipient) == 0 {
			return a, errors.New("at least one recipient must be specified (e.g., --recipient=jwe:mypubkey.pem)")
		}
	}
	// While --recipient can be specified only for `nerdctl image encrypt`,
	// --dec-recipient can be specified for both `nerdctl image encrypt` and `nerdctl image decrypt`.
	a.DecRecipient = options.DecRecipients
	return a, nil
}

// CreateCryptoConfig from the list of recipient strings and list of key paths of private keys
func CreateCryptoConfig(args parsehelpers.EncArgs, descs []ocispec.Descriptor, pubKeys [][]byte) (encconfig.CryptoConfig, error) {
	recipients := args.Recipient
	keys := args.Key

	var decryptCc *encconfig.CryptoConfig
	ccs := []encconfig.CryptoConfig{}
	if len(keys) > 0 {
		dcc, err := parsehelpers.CreateDecryptCryptoConfig(args, descs)
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		decryptCc = &dcc
		ccs = append(ccs, dcc)
	}

	if len(recipients) > 0 {
		encryptCcs := []encconfig.CryptoConfig{}
		if len(pubKeys) > 0 {
			jweCc, err := encconfig.EncryptWithJwe(pubKeys)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			encryptCcs = append(encryptCcs, jweCc)
		}
		ecc := encconfig.CombineCryptoConfigs(encryptCcs)
		if decryptCc != nil {
			ecc.EncryptConfig.AttachDecryptConfig(decryptCc.DecryptConfig)
		}
		ccs = append(ccs, ecc)
	}

	if len(ccs) > 0 {
		return encconfig.CombineCryptoConfigs(ccs), nil
	}
	return encconfig.CryptoConfig{}, nil
}

// CreateCryptoConfig from the list of recipient strings and list of key paths of private keys
func CreateDecryptCryptoConfig(args parsehelpers.EncArgs, descs []ocispec.Descriptor, privKeysPasswords, privKeys [][]byte) (encconfig.CryptoConfig, error) {
	ccs := []encconfig.CryptoConfig{}

	if len(privKeys) > 0 {
		privKeysCc, err := encconfig.DecryptWithPrivKeys(privKeys, privKeysPasswords)
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		ccs = append(ccs, privKeysCc)
	}
	return encconfig.CombineCryptoConfigs(ccs), nil
}
