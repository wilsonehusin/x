package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

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

func (a *Authenticator) EncryptData(secret []byte) ([]byte, error) {
	// NIST: For GCM a 12 byte IV is strongly suggested as other IV lengths will
	// require additional calculations.
	// crypto/cipher: Never use more than 2^32 random nonces with a given key
	// because of the risk of a repeat.
	// TODO: have a plan once 4 billion secrets have been stored?
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return a.authenticator.Seal(nonce, nonce, secret, nil), nil
}

func (a *Authenticator) DecryptData(data []byte) ([]byte, error) {
	nonceSize := a.authenticator.NonceSize()
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	return a.authenticator.Open(nil, nonce, ciphertext, nil)
}
