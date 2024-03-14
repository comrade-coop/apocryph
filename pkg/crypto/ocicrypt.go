package crypto

// ALL methods are copied from
// github.com/containerd/imgcrypt@v1.1.9/images/encryption/parsehelpers/parsehelpers.go
// with a slight adjustment to support taking keys as byte arrays instead of
// reading from the file system

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/images/converter"
	"github.com/containerd/imgcrypt/images/encryption"
	"github.com/containerd/imgcrypt/images/encryption/parsehelpers"
	"github.com/containerd/nerdctl/pkg/api/types"
	"github.com/containerd/nerdctl/pkg/platformutil"
	"github.com/containerd/nerdctl/pkg/referenceutil"
	"github.com/containers/ocicrypt"
	encconfig "github.com/containers/ocicrypt/config"
	"github.com/containers/ocicrypt/config/pkcs11config"
	"github.com/containers/ocicrypt/crypto/pkcs11"
	encutils "github.com/containers/ocicrypt/utils"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// from containerd/nerdctl/pkg/cmd/image/crypt.go modified to take extra keys
// args used for ecryption/decryption
func Crypt(ctx context.Context, client *containerd.Client, srcRawRef, targetRawRef string, encrypt bool, options types.ImageCryptOptions, pubKeys, prvKeys [][]byte) error {
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
		cc, err := CreateDecryptCryptoConfig(imgcryptFlags, layerDescs, pubKeys, prvKeys)
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
func CreateCryptoConfig(args parsehelpers.EncArgs, descs []ocispec.Descriptor, recipientKeys [][]byte) (encconfig.CryptoConfig, error) {
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
		gpgRecipients, pubKeys, x509s, pkcs11Pubkeys, pkcs11Yamls, keyProvider, err := processRecipientKeys(recipients, recipientKeys)
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		encryptCcs := []encconfig.CryptoConfig{}

		gpgClient, err := parsehelpers.CreateGPGClient(args)
		gpgInstalled := err == nil
		if len(gpgRecipients) > 0 && gpgInstalled {
			gpgPubRingFile, err := gpgClient.ReadGPGPubRingFile()
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}

			gpgCc, err := encconfig.EncryptWithGpg(gpgRecipients, gpgPubRingFile)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			encryptCcs = append(encryptCcs, gpgCc)
		}

		// Create Encryption Crypto Config
		if len(x509s) > 0 {
			pkcs7Cc, err := encconfig.EncryptWithPkcs7(x509s)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			encryptCcs = append(encryptCcs, pkcs7Cc)
		}
		if len(pubKeys) > 0 {
			jweCc, err := encconfig.EncryptWithJwe(pubKeys)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			encryptCcs = append(encryptCcs, jweCc)
		}
		var p11conf *pkcs11.Pkcs11Config
		if len(pkcs11Yamls) > 0 || len(pkcs11Pubkeys) > 0 {
			p11conf, err = pkcs11config.GetUserPkcs11Config()
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			pkcs11Cc, err := encconfig.EncryptWithPkcs11(p11conf, pkcs11Pubkeys, pkcs11Yamls)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			encryptCcs = append(encryptCcs, pkcs11Cc)
		}

		if len(keyProvider) > 0 {
			keyProviderCc, err := encconfig.EncryptWithKeyProvider(keyProvider)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			encryptCcs = append(encryptCcs, keyProviderCc)
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
func CreateDecryptCryptoConfig(args parsehelpers.EncArgs, descs []ocispec.Descriptor, pubKeys, prvKeys [][]byte) (encconfig.CryptoConfig, error) {
	ccs := []encconfig.CryptoConfig{}

	// x509 cert is needed for PKCS7 decryption
	_, _, x509s, _, _, _, err := processRecipientKeys(args.DecRecipient, pubKeys)
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}

	gpgSecretKeyRingFiles, gpgSecretKeyPasswords, privKeys, privKeysPasswords, pkcs11Yamls, keyProviders, err := processPrivateKeyFiles(args.Key, prvKeys)
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}

	_, err = parsehelpers.CreateGPGClient(args)
	gpgInstalled := err == nil
	if gpgInstalled {
		if len(gpgSecretKeyRingFiles) == 0 && len(privKeys) == 0 && len(pkcs11Yamls) == 0 && len(keyProviders) == 0 && descs != nil {
			// Get pgp private keys from keyring only if no private key was passed
			gpgPrivKeys, gpgPrivKeyPasswords, err := getGPGPrivateKeys(args, gpgSecretKeyRingFiles, descs, true)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}

			gpgCc, err := encconfig.DecryptWithGpgPrivKeys(gpgPrivKeys, gpgPrivKeyPasswords)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			ccs = append(ccs, gpgCc)

		} else if len(gpgSecretKeyRingFiles) > 0 {
			gpgCc, err := encconfig.DecryptWithGpgPrivKeys(gpgSecretKeyRingFiles, gpgSecretKeyPasswords)
			if err != nil {
				return encconfig.CryptoConfig{}, err
			}
			ccs = append(ccs, gpgCc)

		}
	}

	if len(x509s) > 0 {
		x509sCc, err := encconfig.DecryptWithX509s(x509s)
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		ccs = append(ccs, x509sCc)
	}
	if len(privKeys) > 0 {
		privKeysCc, err := encconfig.DecryptWithPrivKeys(privKeys, privKeysPasswords)
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		ccs = append(ccs, privKeysCc)
	}
	if len(pkcs11Yamls) > 0 {
		p11conf, err := pkcs11config.GetUserPkcs11Config()
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		pkcs11PrivKeysCc, err := encconfig.DecryptWithPkcs11Yaml(p11conf, pkcs11Yamls)
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		ccs = append(ccs, pkcs11PrivKeysCc)
	}
	if len(keyProviders) > 0 {
		keyProviderCc, err := encconfig.DecryptWithKeyProvider(keyProviders)
		if err != nil {
			return encconfig.CryptoConfig{}, err
		}
		ccs = append(ccs, keyProviderCc)
	}
	return encconfig.CombineCryptoConfigs(ccs), nil
}

// Copied from parsehelpers
// processRecipientKeys sorts the array of recipients by type. Recipients may be either
// x509 certificates, public keys, or PGP public keys identified by email address or name
func processRecipientKeys(recipients []string, keys [][]byte) ([][]byte, [][]byte, [][]byte, [][]byte, [][]byte, [][]byte, error) {
	var (
		gpgRecipients [][]byte
		pubkeys       [][]byte
		x509s         [][]byte
		pkcs11Pubkeys [][]byte
		pkcs11Yamls   [][]byte
		keyProvider   [][]byte
	)

	for kid, recipient := range recipients {

		idx := strings.Index(recipient, ":")
		if idx < 0 {
			return nil, nil, nil, nil, nil, nil, errors.New("invalid recipient format")
		}

		protocol := recipient[:idx]
		value := recipient[idx+1:]

		switch protocol {
		case "pgp":
			gpgRecipients = append(gpgRecipients, []byte(value))

		case "jwe":
			if !encutils.IsPublicKey(keys[kid]) {
				return nil, nil, nil, nil, nil, nil, errors.New("file provided is not a public key")
			}
			pubkeys = append(pubkeys, keys[kid])

		case "pkcs7":
			if !encutils.IsCertificate(keys[kid]) {
				return nil, nil, nil, nil, nil, nil, errors.New("file provided is not an x509 cert")
			}
			x509s = append(x509s, keys[kid])

		case "pkcs11":
			if encutils.IsPkcs11PublicKey(keys[kid]) {
				pkcs11Yamls = append(pkcs11Yamls, keys[kid])
			} else if encutils.IsPublicKey(keys[kid]) {
				pkcs11Pubkeys = append(pkcs11Pubkeys, keys[kid])
			} else {
				return nil, nil, nil, nil, nil, nil, errors.New("provided file is not a public key")
			}

		case "provider":
			keyProvider = append(keyProvider, []byte(value))

		default:
			return nil, nil, nil, nil, nil, nil, errors.New("provided protocol not recognized")
		}
	}
	return gpgRecipients, pubkeys, x509s, pkcs11Pubkeys, pkcs11Yamls, keyProvider, nil
}

// processPrivateKeyFiles sorts the different types of private key files; private key files may either be
// private keys or GPG private key ring files. The private key files may include the password for the
// private key and take any of the following forms:
// - <filename>
// - <filename>:file=<passwordfile>
// - <filename>:pass=<password>
// - <filename>:fd=<filedescriptor>
// - <filename>:<password>
// - keyprovider:<...>
func processPrivateKeyFiles(keyFilesAndPwds []string, prvKeys [][]byte) ([][]byte, [][]byte, [][]byte, [][]byte, [][]byte, [][]byte, error) {
	var (
		gpgSecretKeyRingFiles [][]byte
		gpgSecretKeyPasswords [][]byte
		privkeys              [][]byte
		privkeysPasswords     [][]byte
		pkcs11Yamls           [][]byte
		keyProviders          [][]byte
		err                   error
	)
	// keys needed for decryption in case of adding a recipient
	for kid, keyfileAndPwd := range keyFilesAndPwds {
		var password []byte

		// treat "provider" protocol separately
		if strings.HasPrefix(keyfileAndPwd, "provider:") {
			keyProviders = append(keyProviders, []byte(keyfileAndPwd[9:]))
			continue
		}
		parts := strings.Split(keyfileAndPwd, ":")
		if len(parts) == 2 {
			password, err = processPwdString(parts[1])
			if err != nil {
				return nil, nil, nil, nil, nil, nil, err
			}
		}

		keyfile := parts[0]
		isPrivKey, err := encutils.IsPrivateKey(prvKeys[kid], password)
		if encutils.IsPasswordError(err) {
			return nil, nil, nil, nil, nil, nil, err
		}

		if encutils.IsPkcs11PrivateKey(prvKeys[kid]) {
			pkcs11Yamls = append(pkcs11Yamls, privkeys[kid])
		} else if isPrivKey {
			privkeys = append(privkeys, prvKeys[kid])
			privkeysPasswords = append(privkeysPasswords, password)
		} else if encutils.IsGPGPrivateKeyRing(prvKeys[kid]) {
			gpgSecretKeyRingFiles = append(gpgSecretKeyRingFiles, prvKeys[kid])
			gpgSecretKeyPasswords = append(gpgSecretKeyPasswords, password)
		} else {
			return nil, nil, nil, nil, nil, nil, fmt.Errorf("unidentified private key in file %s", keyfile)
		}
	}
	return gpgSecretKeyRingFiles, gpgSecretKeyPasswords, privkeys, privkeysPasswords, pkcs11Yamls, keyProviders, nil
}
func getGPGPrivateKeys(args parsehelpers.EncArgs, gpgSecretKeyRingFiles [][]byte, descs []ocispec.Descriptor, mustFindKey bool) (gpgPrivKeys [][]byte, gpgPrivKeysPwds [][]byte, err error) {
	gpgClient, err := parsehelpers.CreateGPGClient(args)
	if err != nil {
		return nil, nil, err
	}

	var gpgVault ocicrypt.GPGVault
	if len(gpgSecretKeyRingFiles) > 0 {
		gpgVault = ocicrypt.NewGPGVault()
		err = gpgVault.AddSecretKeyRingDataArray(gpgSecretKeyRingFiles)
		if err != nil {
			return nil, nil, err
		}
	}
	return ocicrypt.GPGGetPrivateKey(descs, gpgClient, gpgVault, mustFindKey)
}

// processPwdString process a password that may be in any of the following formats:
// - file=<passwordfile>
// - pass=<password>
// - fd=<filedescriptor>
// - <password>
func processPwdString(pwdString string) ([]byte, error) {
	if strings.HasPrefix(pwdString, "file=") {
		return os.ReadFile(pwdString[5:])
	} else if strings.HasPrefix(pwdString, "pass=") {
		return []byte(pwdString[5:]), nil
	} else if strings.HasPrefix(pwdString, "fd=") {
		fdStr := pwdString[3:]
		fd, err := strconv.Atoi(fdStr)
		if err != nil {
			return nil, fmt.Errorf("could not parse file descriptor %s: %w", fdStr, err)
		}
		f := os.NewFile(uintptr(fd), "pwdfile")
		if f == nil {
			return nil, fmt.Errorf("%s is not a valid file descriptor", fdStr)
		}
		defer f.Close()
		pwd := make([]byte, 64)
		n, err := f.Read(pwd)
		if err != nil {
			return nil, fmt.Errorf("could not read from file descriptor: %w", err)
		}
		return pwd[:n], nil
	}
	return []byte(pwdString), nil
}
