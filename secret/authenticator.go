package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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

// Encrypt takes secret and returns encrypted ciphertext using global authenticator.
func Encrypt(secret []byte) ([]byte, error) {
	if globalAuth == nil {
		return nil, fmt.Errorf("unable to encrypt: global authenticator is not set (use SetGlobal)")
	}
	return globalAuth.Encrypt(secret)
}

// EncryptBase64 is similar to Encrypt, except the output value is now Base64-encoded,
// therefore should be decrypted by DecryptBase64.
func EncryptBase64(secret []byte) ([]byte, error) {
	if globalAuth == nil {
		return nil, fmt.Errorf("unable to encrypt: global authenticator is not set (use SetGlobal)")
	}
	return globalAuth.EncryptBase64(secret)
}

// Decrypt takes ciphertext and returns decrypted secret using global authenticator.
func Decrypt(ciphertext []byte) ([]byte, error) {
	if globalAuth == nil {
		return nil, fmt.Errorf("unable to decrypt: global authenticator is not set (use SetGlobal)")
	}
	return globalAuth.Decrypt(ciphertext)
}

// DecryptBase64 is similar to Decrypt, except it takes input value which was Base64-encoded,
// therefore should only be used for ciphertexts encrypted by EncryptBase64
func DecryptBase64(b64 []byte) ([]byte, error) {
	if globalAuth == nil {
		return nil, fmt.Errorf("unable to decrypt: global authenticator is not set (use SetGlobal)")
	}
	return globalAuth.DecryptBase64(b64)
}

// Encrypt takes in secret and outputs ciphertext
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

// EncryptBase64 is similar to Encrypt, except the output value is now Base64-encoded,
// therefore should be decrypted by DecryptBase64.
func (a *Authenticator) EncryptBase64(secret []byte) ([]byte, error) {
	ciphertext, err := a.Encrypt(secret)
	if err != nil {
		return nil, err
	}
	b64 := make([]byte, base64.RawURLEncoding.EncodedLen(len(ciphertext)))
	base64.RawURLEncoding.Encode(b64, ciphertext)
	return b64, nil
}

// Decrypt takes in ciphertext and outputs secret
func (a *Authenticator) Decrypt(data []byte) ([]byte, error) {
	nonceSize := a.authenticator.NonceSize()
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	return a.authenticator.Open(nil, nonce, ciphertext, nil)
}

// DecryptBase64 is similar to Decrypt, except it takes input value which was Base64-encoded,
// therefore should only be used for ciphertexts encrypted by EncryptBase64
func (a *Authenticator) DecryptBase64(b64 []byte) ([]byte, error) {
	ciphertext := make([]byte, base64.RawURLEncoding.DecodedLen(len(b64)))
	if _, err := base64.RawURLEncoding.Decode(ciphertext, b64); err != nil {
		return nil, err
	}
	return a.Decrypt(ciphertext)
}
