// SPDX-License-Identifier: GPL-3.0

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Helper functions related to AES-GCM encryption, used by the crypto package itself.

const aes_NONCE_SIZE = 12
const aes_KEY_SIZE = 32
const pbkdf2_ITER_COUNT = 10

func encryptWithAES(data []byte, key []byte) ([]byte, error) {
	aesGCM, err := createAESGCMMode(key)
	if err != nil {
		return nil, err
	}
	nonce, err := generateNonce()
	if err != nil {
		return nil, err
	}
	dst := nonce[:]
	ciphertext := aesGCM.Seal(dst, nonce, data, nil)
	return ciphertext, nil
}

func decryptWithAES(data []byte, key []byte) ([]byte, error) {
	aesGCM, err := createAESGCMMode(key)
	if err != nil {
		return nil, err
	}
	nonce := data[:aes_NONCE_SIZE]
	data = data[aes_NONCE_SIZE:]

	plaintext, err := aesGCM.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// Generate a AES_KEY_SIZE-byte random key
func generateAESKey(random io.Reader) ([]byte, error) {
	key := make([]byte, aes_KEY_SIZE)
	if _, err := io.ReadFull(random, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Generate a NONCE_SIZE-byte random nonce
func generateNonce() ([]byte, error) {
	nonce := make([]byte, aes_NONCE_SIZE)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
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
