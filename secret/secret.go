package secret

import (
	"fmt"
)

type Bytes struct {
	authenticator *Authenticator
	secret        []byte
}

func NewBytes(secret []byte) Bytes {
	return Bytes{secret: secret}
}

func NewBytesWithAuth(authenticator *Authenticator, secret []byte) Bytes {
	return Bytes{
		authenticator: authenticator,
		secret:        secret,
	}
}

// MarshalText outputs base64 (URL variant) representation of encrypted secret.
// MarshalJSON was deliberately not added because json.Marshal relies on MarshalText
// for JSON keys.
// RawURLEncoding needs to be used to ensure that results are not padded, otherwise
// it might cause misalignment on encrypt-decrypt process.
// MarshalText will use the attached authenticator if provided, otherwise will
// fallback to globalAuth, configured by SetGlobal
func (s Bytes) MarshalText() ([]byte, error) {
	auth := globalAuth
	if s.authenticator != nil {
		auth = s.authenticator
	}
	if auth == nil {
		return nil, fmt.Errorf("missing authenticator: initialize authenticator or use SetGlobal")
	}
	ciphertext, err := auth.EncryptBase64(s.secret)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// UnmarshalText will use the attached authenticator if provided, otherwise will
// fallback to globalAuth, configured by SetGlobal
func (s *Bytes) UnmarshalText(b64 []byte) error {
	auth := globalAuth
	if s.authenticator != nil {
		auth = s.authenticator
	}
	if auth == nil {
		return fmt.Errorf("missing authenticator: initialize authenticator or use SetGlobal")
	}
	secret, err := auth.DecryptBase64(b64)
	if err != nil {
		return err
	}
	s.secret = secret
	return nil
}

func (s Bytes) MarshalBinary() ([]byte, error) {
	auth := globalAuth
	if s.authenticator != nil {
		auth = s.authenticator
	}
	if auth == nil {
		return nil, fmt.Errorf("missing authenticator: initialize authenticator or use SetGlobal")
	}
	ciphertext, err := auth.Encrypt(s.secret)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

func (s *Bytes) UnmarshalBinary(b []byte) error {
	auth := globalAuth
	if s.authenticator != nil {
		auth = s.authenticator
	}
	if auth == nil {
		return fmt.Errorf("missing authenticator: initialize authenticator or use SetGlobal")
	}
	secret, err := auth.Decrypt(b)
	if err != nil {
		return err
	}
	s.secret = secret
	return nil
}

func (s Bytes) SetValue(b []byte) {
	s.secret = b
}

func (s Bytes) Value() []byte {
	return s.secret
}

type String struct {
	Bytes
}

func NewString(secret string) String {
	return String{
		Bytes: NewBytes([]byte(secret)),
	}
}

func NewStringWithAuth(authenticator *Authenticator, secret string) String {
	return String{
		Bytes: NewBytesWithAuth(authenticator, []byte(secret)),
	}
}

func (s String) SetValue(str string) {
	s.secret = []byte(str)
}

func (s String) Value() string {
	return string(s.secret)
}
