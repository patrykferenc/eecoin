package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

var (
	ErrPrivateKeyNotFound = errors.New("private key not found")
	NoKeysFound           = errors.New("no keys found")
)

const defaultRsaBitSize = 2048

type RsaWallet struct {
	mainId *Key[*rsa.PrivateKey, crypto.PublicKey]
	keys   map[crypto.PublicKey]privateKeyElement
}

type privateKeyElement struct {
	key     *rsa.PrivateKey
	present bool
}

func (w *RsaWallet) SetMainIdentity(key *Key[*rsa.PrivateKey, crypto.PublicKey]) error {
	if key.private != nil {
		w.mainId = &Key[*rsa.PrivateKey, crypto.PublicKey]{private: key.private, public: key.private.PublicKey, algType: RSA}
		return nil
	}
	return ErrPrivateKeyNotFound
}

func (w *RsaWallet) Add(key Key[*rsa.PrivateKey, crypto.PublicKey]) error {
	if key.private != nil {
		pub := key.private.Public()
		w.keys[pub] = privateKeyElement{key: key.private, present: true}
		return nil
	} else if key.public != nil {
		w.keys[key.public] = privateKeyElement{key: nil, present: false}
	}
	return NoKeysFound
}
func (w *RsaWallet) Type() Algorithm {
	return RSA
}

func NewRsaKey() (Key[*rsa.PrivateKey, crypto.PublicKey], error) {
	key, err := rsa.GenerateKey(rand.Reader, defaultRsaBitSize)
	if err != nil {
		return Key[*rsa.PrivateKey, crypto.PublicKey]{algType: RSA}, nil
	}
	return Key[*rsa.PrivateKey, crypto.PublicKey]{private: key, public: key.Public(), algType: RSA}, nil
}

func NewRsaWallet(mainId *Key[*rsa.PrivateKey, crypto.PublicKey]) *RsaWallet {
	wallet := &RsaWallet{keys: make(map[crypto.PublicKey]privateKeyElement), mainId: mainId}
	_ = wallet.Add(*mainId)
	return wallet
}
