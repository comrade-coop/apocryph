package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	encconfig "github.com/containers/ocicrypt/config"
	"github.com/go-jose/go-jose/v3"
)

const (
	KeyTypeEcdsaP256 = "ECDSAP256"
	KeyTypeRsa4096   = "RSA4096"
	KeyTypeAESGCM256 = "AESGCM256"
)

var KeyTypeOcicrypt = KeyTypeRsa4096 // NOTE: ocicrypt JWE keywrap does not allow us to use any JWK key algorithms other than RSA_OAEP
var KeyTypeEncrypt = KeyTypeAESGCM256

func CreateKey(keys *[]*pb.Key, keyType string) (int32, error) {
	var err error
	var key interface{}
	var keyAlg string
	switch keyType {
	case KeyTypeEcdsaP256:
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		keyAlg = string(jose.ECDH_ES_A256KW)
	case KeyTypeRsa4096:
		key, err = rsa.GenerateKey(rand.Reader, 4096)
		keyAlg = string(jose.RSA_OAEP)
	case KeyTypeAESGCM256:
		key, err = GenerateAESKey(rand.Reader)
		keyAlg = string(jose.DIRECT)
	}
	if err != nil {
		return -1, err
	}
	keyIdx := int32(len(*keys))
	jwk, err := jose.JSONWebKey{
		Key:       key,
		KeyID:     fmt.Sprintf("key-%d", keyIdx),
		Algorithm: keyAlg,
	}.MarshalJSON()
	if err != nil {
		return -1, err
	}
	*keys = append(*keys, &pb.Key{
		Data: jwk,
	})
	return keyIdx, nil
}

func getJwk(keys []*pb.Key, keyIdx int32) (*jose.JSONWebKey, error) {
	if keyIdx < 0 {
		return nil, nil
	}
	if int(keyIdx) >= len(keys) {
		return nil, errors.New("Invalid keyIdx")
	}
	jwk := keys[keyIdx].Data
	jwkUnmarshalled := &jose.JSONWebKey{}
	err := jwkUnmarshalled.UnmarshalJSON(jwk)
	if err != nil {
		return nil, err
	}
	return jwkUnmarshalled, nil
}

func getAESKey(keys []*pb.Key, keyIdx int32) ([]byte, error) {
	jwkUnmarshalled, err := getJwk(keys, keyIdx)
	if err != nil {
		return nil, err
	}
	if jwkUnmarshalled == nil {
		return nil, nil
	}
	switch key := jwkUnmarshalled.Key.(type) {
	case []byte:
		return key, nil
	default:
		return nil, errors.New("Not an AES key")
	}
}

func EncryptWith(keys []*pb.Key, keyIdx int32, plaintext []byte) ([]byte, error) {
	aesKey, err := getAESKey(keys, keyIdx)
	if err != nil {
		return nil, err
	}
	if aesKey == nil {
		return plaintext, nil
	}
	return EncryptWithAES(plaintext, aesKey)
}

func DecryptWith(keys []*pb.Key, keyIdx int32, cyphertext []byte) ([]byte, error) {
	aesKey, err := getAESKey(keys, keyIdx)
	if err != nil {
		return nil, err
	}
	if aesKey == nil {
		return cyphertext, nil
	}
	return DecryptWithAES(cyphertext, aesKey)
}

func GetCryptoConfig(keys []*pb.Key, keyIdx int32) (encconfig.CryptoConfig, error) {
	jwkUnmarshalled, err := getJwk(keys, keyIdx)
	if err != nil || jwkUnmarshalled == nil {
		return encconfig.CryptoConfig{}, err
	}
	jwk, err := jwkUnmarshalled.MarshalJSON()
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}
	jwkPub, err := jwkUnmarshalled.Public().MarshalJSON()
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}
	encryptConfig, err := encconfig.EncryptWithJwe([][]byte{jwkPub})
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}
	decryptConfig, err := encconfig.DecryptWithPrivKeys([][]byte{jwk}, [][]byte{nil})
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}
	return encconfig.CombineCryptoConfigs([]encconfig.CryptoConfig{
		encryptConfig,
		decryptConfig,
	}), nil
}
