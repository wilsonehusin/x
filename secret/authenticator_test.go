package secret

import "testing"

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

func TestHMACEqual(t *testing.T) {
	msg := []byte(`carby is best kirby`)

	auth := getAuth()

	hmac, err := auth.HMAC(msg)
	if err != nil {
		t.Fatal(err)
	}
	if hmac == nil || len(hmac) == 0 {
		t.Fatal("hmac was not calculated")
	}

	if err := auth.HMACCheck(msg, hmac); err != nil {
		t.Fatal(err)
	}
}

func TestHMACNotEqual(t *testing.T) {
	original := []byte(`mario`)
	counterfeit := []byte(`luigi`)

	auth := getAuth()

	hmac, err := auth.HMAC(original)
	if err != nil {
		t.Fatal(err)
	}
	if hmac == nil || len(hmac) == 0 {
		t.Fatal("hmac was not calculated")
	}

	err = auth.HMACCheck(counterfeit, hmac)
	if err != ErrHMACMismatch {
		if err == nil {
			t.Fatal("expecting ErrHMACMismatch, but received nil")
		} else {
			t.Fatalf("expecting ErrHMACMismatch, but received %s", err.Error())
		}
	}
}
