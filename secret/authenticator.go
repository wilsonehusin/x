package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

var globalAuth *Authenticator

type Authenticator struct {
	authenticator cipher.AEAD
}

func NewAuthenticatorAESGCM(key []byte) (*Authenticator, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Authenticator{authenticator: aead}, nil
}

func SetGlobal(a *Authenticator) {
	globalAuth = a
}

func Encrypt(secret []byte) ([]byte, error) {
	if globalAuth == nil {
		return nil, fmt.Errorf("unable to encrypt: global authenticator is not set (use SetGlobal)")
	}
	return globalAuth.Encrypt(secret)
}

func Decrypt(ciphertext []byte) ([]byte, error) {
	if globalAuth == nil {
		return nil, fmt.Errorf("unable to decrypt: global authenticator is not set (use SetGlobal)")
	}
	return globalAuth.Decrypt(ciphertext)
}

func (a *Authenticator) Encrypt(secret []byte) ([]byte, error) {
	// NIST: For GCM a 12 byte IV is strongly suggested as other IV lengths will
	// require additional calculations.
	// crypto/cipher: Never use more than 2^32 random nonces with a given key
	// because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return a.authenticator.Seal(nonce, nonce, secret, nil), nil
}

func (a *Authenticator) Decrypt(data []byte) ([]byte, error) {
	nonceSize := a.authenticator.NonceSize()
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	return a.authenticator.Open(nil, nonce, ciphertext, nil)
}
