package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"

	pb "github.com/comrade-coop/trusted-pods/pkg/proto"
	encconfig "github.com/containers/ocicrypt/config"
	"github.com/go-jose/go-jose/v3"
)

const (
	KeyTypeEcdsaP256 = "ECDSAP256"
	KeyTypeRsa4096   = "RSA4096"
	KeyTypeAESGCM256 = "AESGCM256"
)

var KeyTypeOcicrypt = KeyTypeEcdsaP256
var KeyTypeEncrypt = KeyTypeAESGCM256

func NewKey(keyType string) (*pb.Key, error) {
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
		return nil, err
	}
	jwk, err := jose.JSONWebKey{
		Key:       key,
		KeyID:     "key",
		Algorithm: keyAlg,
	}.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &pb.Key{
		Data: jwk,
	}, nil
}

func getJwkKey(key *pb.Key) (*jose.JSONWebKey, error) {
	jwk := key.Data
	jwkUnmarshalled := &jose.JSONWebKey{}
	err := jwkUnmarshalled.UnmarshalJSON(jwk)
	if err != nil {
		return nil, err
	}
	return jwkUnmarshalled, nil
}

func getAESKey(key *pb.Key) ([]byte, error) {
	jwkUnmarshalled, err := getJwkKey(key)
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

func EncryptWithKey(key *pb.Key, plaintext []byte) ([]byte, error) {
	aesKey, err := getAESKey(key)
	if err != nil {
		return nil, err
	}
	if aesKey == nil {
		return plaintext, nil
	}
	return EncryptWithAES(plaintext, aesKey)
}

func DecryptWithKey(key *pb.Key, cyphertext []byte) ([]byte, error) {
	aesKey, err := getAESKey(key)
	if err != nil {
		return nil, err
	}
	if aesKey == nil {
		return cyphertext, nil
	}
	return DecryptWithAES(cyphertext, aesKey)
}

func GetCryptoConfigKey(key *pb.Key) (encconfig.CryptoConfig, error) {
	jwkUnmarshalled, err := getJwkKey(key)
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}
	jwk, err := jwkUnmarshalled.MarshalJSON()
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}
	keyPub, err := x509.MarshalPKIXPublicKey(jwkUnmarshalled.Public().Key) // TODO: https://github.com/containers/ocicrypt/pull/99 - jwkUnmarshalled.Public().MarshalJSON()
	if err != nil {
		return encconfig.CryptoConfig{}, err
	}
	encryptConfig, err := encconfig.EncryptWithJwe([][]byte{keyPub})
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
