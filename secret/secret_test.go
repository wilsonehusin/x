package secret

import (
	"bytes"
	"encoding/json"
	"testing"
)

const (
	testStringKey = "955880d5f4f43c66751848c06fedb78e420995b373418dcfb856ca559deb71c3"
)

func getAuth() *Authenticator {
	key, err := KeyFromString(testStringKey)
	if err != nil {
		panic(err)
	}
	auth, err := NewAuthenticatorAESGCM(key)
	if err != nil {
		panic(err)
	}
	return auth
}

func TestBytesEncodeDecodeJSON(t *testing.T) {
	auth := getAuth()

	lines := []string{
		`we're no strangers to love`,
		`you know the rules and so do i`,
	}

	for _, line := range lines {
		src := NewBytes(auth, []byte(line))
		raw, err := json.Marshal(src)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("ciphertext (json string): %s", raw)
		dst := NewBytes(auth, nil)
		if err := json.Unmarshal(raw, dst); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(src.Value(), dst.Value()) {
			t.Fatalf("unequal:\n\tsrc: %x\n\tdst: %x\n", src.Value(), dst.Value())
		}
	}
}

func TestStringEncodeDecodeJSON(t *testing.T) {
	auth := getAuth()

	lines := []string{
		`a full commitment's what i'm thinking of`,
		`you wouldn't get this from any other guy`,
	}
	for _, line := range lines {
		src := NewString(auth, line)
		raw, err := json.Marshal(src)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("ciphertext (json string): %s", raw)
		dst := NewString(auth, "")
		if err := json.Unmarshal(raw, dst); err != nil {
			t.Fatal(err)
		}
		if src.Value() != dst.Value() {
			t.Fatalf("unequal:\n\tsrc: %s\n\tdst:%s\n", src.Value(), dst.Value())
		}
	}
}

func TestStructWithStringEncodeDecodeJSON(t *testing.T) {
	auth := getAuth()

	type clientConfig struct {
		ClientID, ClientSecret *String
	}

	src := &clientConfig{
		ClientID:     NewString(auth, "this-is-client-id"),
		ClientSecret: NewString(auth, "this-is-client-secret"),
	}

	raw, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ciphertext (json object): %s", raw)

	dst := &clientConfig{
		ClientID:     NewString(auth, ""),
		ClientSecret: NewString(auth, ""),
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		t.Fatal(err)
	}
	if src.ClientID.Value() != dst.ClientID.Value() {
		t.Fatalf("unequal:\n\tsrc: %+v\n\tdst: %+v\n", src.ClientID.Value(), dst.ClientID.Value())
	}
	if src.ClientSecret.Value() != dst.ClientSecret.Value() {
		t.Fatalf("unequal:\n\tsrc: %+v\n\tdst: %+v\n", src.ClientSecret.Value(), dst.ClientSecret.Value())
	}
}
