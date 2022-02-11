# `go.husin.dev/x/secret`

[![Go Reference](https://pkg.go.dev/badge/go.husin.dev/x/secret.svg)](https://pkg.go.dev/go.husin.dev/x/secret)

Automatic secret encryption, inspired by [`:attr_encrypted`](https://github.com/attr-encrypted/attr_encrypted) in Rails.

Please be advised that this project has not been audited for security, thus use your own discretion.

---

Typically, an application may have a `struct` which contains sensitive information which you might wish to be encrypted at rest. 

Let's take the following example:

```go
type OAuthConfig struct {
  ClientID      string
  ClientSecret  string
}
```

When the above is converted to JSON to be saved, the output will be literal plain-text, such as:

```json
{
  "ClientID": "this-is-client-id",
  "ClientSecret": "this-is-client-secret"
}
```

By using this library, one can define their `struct` as:

```go
type OAuthConfig struct {
  ClientID      *secret.String
  ClientSecret  *secret.String
}
```

And upon being converted to JSON, the output will be encrypted, like:

```json
{
  "ClientID": "0AiBPBwxvs1VB90-OVJP_WA4xs8VUntmgrTdGx_9_ix8ReQgtU99NP9ZMLwK",
  "ClientSecret": "qkUP-UdTT1p6zQjk23aHEFjzoRR2T4j0S4Vxah1wqZIkJKI5-um1czKLYiXE2-YKUQ"
}
```

Meanwhile in the rest of code, those values can be retrieved as such:

```go
var clientID secret.Secret[string]
clientID = secret.NewString("this-is-client-id")
fmt.Println(clientID.Value())

// Out:
// this-is-client-id
```

## Caveats

### Unmarshaling

Upon unmarshaling, `Authentication` must already be known, i.e.:

```go
dst := OAuthConfig{ClientID: NewString(auth, ""), ClientSecret: NewString(auth, "")}
if err := json.Unmarshal(rawData, &dst); err != nil {
  panic(err)
}
```

### Nonce (or "why Marshal() calls are not idempotent?")

Under the hood, the encryption being performed is AES-GCM, provided by Go standard library.
Encryption of data in AES-GCM protocol requires: key, nonce, secret.

Key is the encryption key which should be stored safely.
Nonce is generated at runtime to achieve the goal of salting / IV (initialization vector).

Each call to `MarshalJSON()` generates a new nonce and therefore generates different ciphertext, though all of them are can be decrypted with the same key, because nonce is part of the ciphertext.
