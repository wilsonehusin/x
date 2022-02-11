package secret

import (
	"encoding/base64"
)

type Bytes struct {
	authenticator *Authenticator
	secret        []byte
}

func NewBytes(authenticator *Authenticator, secret []byte) *Bytes {
	return &Bytes{
		authenticator: authenticator,
		secret:        secret,
	}
}

// MarshalText outputs base64 (URL variant) representation of encrypted secret.
// MarshalJSON was deliberately not added because json.Marshal relies on MarshalText
// for JSON keys.
// RawURLEncoding needs to be used to ensure that results are not padded, otherwise
// it might cause misalignment on encrypt-decrypt process.
func (s *Bytes) MarshalText() ([]byte, error) {
	ciphertext, err := s.authenticator.EncryptData(s.secret)
	if err != nil {
		return nil, err
	}
	b64 := make([]byte, base64.RawURLEncoding.EncodedLen(len(ciphertext)))
	base64.RawURLEncoding.Encode(b64, ciphertext)
	return b64, nil
}

func (s *Bytes) UnmarshalText(b64 []byte) error {
	ciphertext := make([]byte, base64.RawURLEncoding.DecodedLen(len(b64)))
	if _, err := base64.RawURLEncoding.Decode(ciphertext, b64); err != nil {
		return err
	}
	secret, err := s.authenticator.DecryptData(ciphertext)
	if err != nil {
		return err
	}
	s.secret = secret
	return nil
}

func (s *Bytes) SetValue(b []byte) {
	s.secret = b
}

func (s *Bytes) Value() []byte {
	return s.secret
}

type String struct {
	*Bytes
}

func NewString(authenticator *Authenticator, secret string) *String {
	return &String{
		Bytes: NewBytes(authenticator, []byte(secret)),
	}
}

func (s *String) SetValue(str string) {
	s.secret = []byte(str)
}

func (s *String) Value() string {
	return string(s.secret)
}
