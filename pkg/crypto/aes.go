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

func AESEncrypt(data []byte, key []byte) ([]byte, error) {
	aesGCM, err := createAESGCMMode(key)
	if err != nil {
		return nil, err
	}
	nonce, err := CreateRandomNonce()
	if err != nil {
		return nil, err
	}
	dst := nonce[:]
	ciphertext := aesGCM.Seal(dst, nonce, data, nil)
	return ciphertext, nil
}

func AESDecrypt(data []byte, key []byte) ([]byte, error) {
	aesGCM, err := createAESGCMMode(key)
	if err != nil {
		return nil, err
	}
	nonce := data[:NONCE_SIZE]
	data = data[NONCE_SIZE:]

	plaintext, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// Generate a AES_KEY_SIZE-byte random key
func CreateRandomKey() ([]byte, error) {
	key := make([]byte, AES_KEY_SIZE)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Generate a NONCE_SIZE-byte random nonce
func CreateRandomNonce() ([]byte, error) {
	nonce := make([]byte, NONCE_SIZE)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// could drive a nonce or a symmetric key, both are []byte
func DeriveKey(password []byte, salt []byte) []byte {
	return pbkdf2.Key(password, salt, ITER_COUNT, AES_KEY_SIZE, sha256.New)
}

func createAESGCMMode(key []byte) (cipher.AEAD, error) {
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
