package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	ErrPrivateKeyNotFound = errors.New("private key not found")
	NoKeysFound           = errors.New("no keys found")
	PemParseError         = errors.New("pem parse error")
)

const defaultRsaBitSize = 2048

type Rsa struct {
	mainId *Key[*rsa.PrivateKey, crypto.PublicKey]
	keys   map[crypto.PublicKey]privateKeyElement
}

type privateKeyElement struct {
	key     *rsa.PrivateKey
	present bool
}

func (w *Rsa) SetMainIdentity(key *Key[*rsa.PrivateKey, crypto.PublicKey]) error {
	if key.private != nil {
		w.mainId = &Key[*rsa.PrivateKey, crypto.PublicKey]{private: key.private, public: key.private.PublicKey, algType: RSA}
		return nil
	}
	return ErrPrivateKeyNotFound
}

func (w *Rsa) Add(key Key[*rsa.PrivateKey, crypto.PublicKey]) error {
	if key.private != nil {
		pub := key.private.Public()
		w.keys[pub] = privateKeyElement{key: key.private, present: true}
		return nil
	} else if key.public != nil {
		w.keys[key.public] = privateKeyElement{key: nil, present: false}
	}
	return NoKeysFound
}
func (w *Rsa) Type() Algorithm {
	return RSA
}

func NewRsaKey() (Key[*rsa.PrivateKey, crypto.PublicKey], error) {
	key, err := rsa.GenerateKey(rand.Reader, defaultRsaBitSize)
	if err != nil {
		return Key[*rsa.PrivateKey, crypto.PublicKey]{algType: RSA}, nil
	}
	return Key[*rsa.PrivateKey, crypto.PublicKey]{private: key, public: key.Public(), algType: RSA}, nil
}

func NewRsaWallet(mainId *Key[*rsa.PrivateKey, crypto.PublicKey]) *Rsa {
	wallet := &Rsa{keys: make(map[crypto.PublicKey]privateKeyElement), mainId: mainId}
	_ = wallet.Add(*mainId)
	return wallet
}

func FromPem(pemData []byte) (Key[*rsa.PrivateKey, crypto.PublicKey], error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return Key[*rsa.PrivateKey, crypto.PublicKey]{}, PemParseError
	}
	key, _ := x509.ParsePKCS1PublicKey(rest)
	return Key[*rsa.PrivateKey, crypto.PublicKey]{public: key, algType: RSA}, nil

}

func ToPem(key Key[*rsa.PrivateKey, crypto.PublicKey]) []byte {
	pub := key.private.Public().(*rsa.PublicKey)
	block := x509.MarshalPKCS1PublicKey(pub)
	return block
}
