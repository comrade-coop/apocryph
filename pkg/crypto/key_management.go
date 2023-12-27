// SPDX-License-Identifier: GPL-3.0

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"

	pb "github.com/comrade-coop/apocryph/pkg/proto"
	encconfig "github.com/containers/ocicrypt/config"
	"github.com/go-jose/go-jose/v3"
)

type KeyType string

const (
	KeyTypeEcdsaP256 = KeyType("ECDSAP256")
	KeyTypeRsa4096   = KeyType("RSA4096")
	KeyTypeAESGCM256 = KeyType("AESGCM256")
)

// Key type to use for Ocicrypt (image encryption) operations.
var KeyTypeOcicrypt = KeyTypeEcdsaP256

// Key type to use for general encryption (secret encryption) operations.
var KeyTypeEncrypt = KeyTypeAESGCM256

// Generate a new [pb.Key] of the given type. Note that for public/private cryptography, this always return a private key.
func NewKey(keyType KeyType) (*pb.Key, error) {
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
		key, err = generateAESKey(rand.Reader)
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

// Get the [jose.JSONWebKey] out of a [pb.Key].
func getJwkKey(key *pb.Key) (*jose.JSONWebKey, error) {
	jwk := key.Data
	jwkUnmarshalled := &jose.JSONWebKey{}
	err := jwkUnmarshalled.UnmarshalJSON(jwk)
	if err != nil {
		return nil, err
	}
	return jwkUnmarshalled, nil
}

// Get an AES key out of a [pb.Key].
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

// Encrypt a given plaintext with a key. Might not work with keys that are not of type [KeyTypeEncrypt].
func EncryptWithKey(key *pb.Key, plaintext []byte) ([]byte, error) {
	aesKey, err := getAESKey(key)
	if err != nil {
		return nil, err
	}
	if aesKey == nil {
		return plaintext, nil
	}
	return encryptWithAES(plaintext, aesKey)
}

// Decrypt a given ciphertext with a key. Might not work with keys that are not of type [KeyTypeEncrypt].
func DecryptWithKey(key *pb.Key, ciphertext []byte) ([]byte, error) {
	aesKey, err := getAESKey(key)
	if err != nil {
		return nil, err
	}
	if aesKey == nil {
		return ciphertext, nil
	}
	return decryptWithAES(ciphertext, aesKey)
}

// Get an OCI [encconfig.CryptoConfig] for a given key. Might not work if the key is not of type [KeyTypeOcicrypt].
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
