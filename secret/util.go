package secret

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	AES128KeyLength = 16 // AES-128 requires keys to be 16-bytes, per "crypto/cipher" documentation
	AES256KeyLength = 32 // AES-256 requires keys to be 16-bytes, per "crypto/cipher" documentation
)

func NewKey(byteLength int) ([]byte, error) {
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func NewStringKey(byteLength int) (string, error) {
	b, err := NewKey(byteLength)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func KeyFromString(str string) ([]byte, error) {
	return hex.DecodeString(str)
}
