package dns_certs

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/dchenk/mazewire/pkg/util"
	"github.com/xenolf/lego/acme"
)

// An Account represents a user's credentials; implements acme.User.
type Account struct {
	Email        string                     `json:"email"`
	Registration *acme.RegistrationResource `json:"registration"`
	key          crypto.PrivateKey
}

// GetEmail returns the email address for the account.
func (a *Account) GetEmail() string { return a.Email }

// GetRegistration returns the server registration
func (a *Account) GetRegistration() *acme.RegistrationResource { return a.Registration }

// GetPrivateKey returns the private account key.
func (a *Account) GetPrivateKey() crypto.PrivateKey { return a.key }

// genPrivateKey generates an ECDSA private key for the user, updating the Account struct's key
// field and saves the PEM-encoded key file.
func (a *Account) genPrivateKey() error {

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}

	pemKey := &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}

	wr := gcsBucket.Object(sslUserKeyFile).NewWriter(context.Background())
	wr.ContentType = util.ContentTypeTextPlain

	pem.Encode(wr, pemKey)

	a.key = privateKey

	return wr.Close()

}

func (a *Account) save() error {
	jsonBytes, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	wr := gcsBucket.Object(sslUserFile).NewWriter(context.Background())
	wr.ContentType = util.ContentTypeTextPlain
	if _, err := wr.Write(jsonBytes); err != nil {
		return err
	}
	return wr.Close()
}

func loadPrivateKey() (crypto.PrivateKey, error) {

	rdr, err := gcsBucket.Object(sslUserKeyFile).NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	keyBytes, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyBytes)

	switch keyBlock.Type {
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(keyBlock.Bytes)
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	}

	return nil, errors.New("unknown private key type")

}
