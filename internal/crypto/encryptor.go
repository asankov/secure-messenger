package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

// Encryptor is the struct that is responsible
// for encrypting and decrypting messages.
type Encryptor struct {
	aead cipher.AEAD
}

// NewEncryptor creates a new Encryptor with the given secret key.
//
// The secret key must be 16, 24 or 32-byte sized.
func NewEncryptor(key string) (*Encryptor, error) {
	if l := len(key); l != 16 && l != 24 && l != 32 {
		return nil, fmt.Errorf("key is of wrong size [%d], required size is either 16, 24 or 32 bytes", l)
	}

	// create an AES Cipher with the secret key
	aes, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("error while creating an AES cipher: %w", err)
	}

	// create a GCM Cipher mode from the AES Cipher
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, fmt.Errorf("error while creating GCM cipher mode: %w", err)
	}

	return &Encryptor{
		aead: gcm,
	}, nil
}

// Encrypt encrypt a plaintext message into a ciphertext.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	// Randomly generate a nonce with the needed size.
	nonce := make([]byte, e.aead.NonceSize())
	n, err := rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("error while generating nonce: %w", err)
	} else if n != len(nonce) {
		// according to the documentation of rand.Read we should never get here,
		// because if n != len(b) then err will be non-nil,
		// but let's do this check just to be safe.
		return "", errors.New("unable to generate nonce with the desired size")
	}

	// Encrypt the text.
	// The encrypted result will be nonce+ciphertext.
	// That way, during decryption, just by knowing the nonce size
	// we can separate it from the ciphertext.
	ciphertext := e.aead.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode the ciphertext to base64 so that we can
	// more easily work with the result as strings
	// (write it to disk, send it over the network, etc.)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a ciphertext back into plaintext.
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	dec, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("error while base64-decoding the ciphertext: %w", err)
	}
	ciphertext = string(dec)

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := e.aead.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := e.aead.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", fmt.Errorf("error while decrypting the ciphertext: %w", err)
	}

	return string(plaintext), nil
}
