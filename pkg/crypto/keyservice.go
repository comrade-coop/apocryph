package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const NONCE_SIZE = 12
const AES_KEY_SIZE = 32
const ITER_COUNT = 10

func AESEncryptWith(data []byte, key []byte, nonce ...[]byte) ([]byte, []byte, error) {
	aesGCM, err := CreateAESGCMMode(key)
	if err != nil {
		return nil, nil, err
	}
	var n []byte
	if len(nonce) > 0 {
		n = nonce[0]
	} else {
		n, err = CreateRandomNonce()
		if err != nil {
			return nil, nil, err
		}

	}
	ciphertext := aesGCM.Seal(nil, n, data, nil)
	return ciphertext, n, nil
}
func AESDecryptWith(data []byte, key []byte, nonce []byte) ([]byte, error) {

	aesGCM, err := CreateAESGCMMode(key)
	// Decrypt the ciphertext using the same key and nonce
	plaintext, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
func CreateRandomKey() ([]byte, error) {
	return createRandomKey()
}

func CreateRandomNonce() ([]byte, error) {
	return createRandomNonce()
}

// could drive a nonce or a symmetric key, both are []byte
func DeriveKey(password, salt []byte, keysize int) []byte {
	return pbkdf2.Key(password, salt, ITER_COUNT, keysize, sha256.New)
}

func createRandomKey() ([]byte, error) {

	// Generate a 32-byte random key
	key := make([]byte, AES_KEY_SIZE)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func createRandomNonce() ([]byte, error) {
	// Generate a 12-byte random nonce
	nonce := make([]byte, NONCE_SIZE)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil

}

func CreateAESGCMMode(key []byte) (cipher.AEAD, error) {

	// Create a cipher block from the key
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Create a GCM with the block
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesGCM, nil

}

type KeyNoncePair struct {
	Key   []byte
	Nonce []byte
}

func CreateKeyNoncePair(password, salt, noncesalt []byte) (*KeyNoncePair, error) {
	key := DeriveKey(password, salt, AES_KEY_SIZE)
	nonce := DeriveKey(password, noncesalt, NONCE_SIZE)
	return &KeyNoncePair{
		Key:   key,
		Nonce: nonce,
	}, nil
}
