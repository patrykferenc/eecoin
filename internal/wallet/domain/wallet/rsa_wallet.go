package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/teris-io/shortid"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrPrivateKeyNotFound = errors.New("private Key not found")
	NoKeysFound           = errors.New("no keys found")
	PemParseError         = errors.New("pem parse error")
)

const defaultRsaBitSize = 2048

type Rsa struct {
	MainId *Key[*rsa.PrivateKey, crypto.PublicKey]
	Keys   map[crypto.PublicKey]KeyElement
}

type KeyElement struct {
	Key     *rsa.PrivateKey
	Present bool
}

func (w *Rsa) SetMainIdentity(key *Key[*rsa.PrivateKey, crypto.PublicKey]) error {
	if key.private != nil {
		w.MainId = &Key[*rsa.PrivateKey, crypto.PublicKey]{private: key.private, Public: key.private.PublicKey, algType: RSA}
		return nil
	}
	return ErrPrivateKeyNotFound
}

func (w *Rsa) Add(key Key[*rsa.PrivateKey, crypto.PublicKey]) error {
	if key.private != nil {
		pub := key.private.Public()
		w.Keys[pub] = KeyElement{Key: key.private, Present: true}
		slog.Info("private-Public Key added to wallet", "pub", pub)
		return nil
	} else if key.Public != nil {
		w.Keys[key.Public] = KeyElement{Key: nil, Present: false}
		slog.Info("Public Key added to wallet", "pub", key.Public)
		return nil
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
	return Key[*rsa.PrivateKey, crypto.PublicKey]{private: key, Public: key.Public(), algType: RSA}, nil
}

func NewRsaWallet(mainId *Key[*rsa.PrivateKey, crypto.PublicKey]) *Rsa {
	wallet := &Rsa{Keys: make(map[crypto.PublicKey]KeyElement), MainId: mainId}
	_ = wallet.Add(*mainId)
	return wallet
}

func newRsaWalletWithoutId() *Rsa {
	wallet := &Rsa{Keys: make(map[crypto.PublicKey]KeyElement)}
	return wallet
}

func PublicFromPem(pemData []byte) (Key[*rsa.PrivateKey, crypto.PublicKey], error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return Key[*rsa.PrivateKey, crypto.PublicKey]{}, PemParseError
	}
	key, e := x509.ParsePKCS1PublicKey(rest)
	if e != nil {
		slog.Error(e.Error())
		return Key[*rsa.PrivateKey, crypto.PublicKey]{}, PemParseError
	}
	return Key[*rsa.PrivateKey, crypto.PublicKey]{Public: key, algType: RSA}, nil
}

func PrivateFromPem(pemData []byte) (Key[*rsa.PrivateKey, crypto.PublicKey], error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return Key[*rsa.PrivateKey, crypto.PublicKey]{}, PemParseError
	}
	key, _ := x509.ParsePKCS1PrivateKey(rest)
	return Key[*rsa.PrivateKey, crypto.PublicKey]{private: key, Public: key.Public(), algType: RSA}, nil
}

func PrivateToPem(key Key[*rsa.PrivateKey, crypto.PublicKey]) []byte {
	priv := key.private
	block := x509.MarshalPKCS1PrivateKey(priv)
	return block
}

func PublicToPem(key Key[*rsa.PrivateKey, crypto.PublicKey]) []byte {
	pub := key.Public.(*rsa.PublicKey)
	block := x509.MarshalPKCS1PublicKey(pub)
	return block
}

func ReadWalletFromDirectory(path string, passphrase *string) (*Rsa, error) {

	absolutePath := filepath.Clean(path)
	directory, err := os.ReadDir(absolutePath)
	wallet := newRsaWalletWithoutId()

	if err == nil {
		for _, file := range directory {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".pub") {
				slog.Info("Opened path", "path", file.Name())
				importPublicKey(absolutePath, file, wallet)
			} else if !file.IsDir() && strings.HasSuffix(file.Name(), ".priv") {
				slog.Info("Opened path", "path", file.Name())
				importMainIdentityPrivateKey(absolutePath, file, passphrase, wallet)
			}
		}
		return wallet, nil
	}

	return nil, err
}

func (w *Rsa) ExportWallet(path string, globalPassphrase *string) error {
	mainId := *w.MainId
	mainIdPem := PrivateToPem(mainId)
	_ = SaveToDirectory(path, "main.priv", mainIdPem, globalPassphrase)
	slog.Info("Exported main identity Key", "path", filepath.Join(path, "main.priv"))

	for pub, priv := range w.Keys {
		if priv.Present && priv.Key != mainId.private {
			key := Key[*rsa.PrivateKey, crypto.PublicKey]{private: priv.Key, algType: RSA}
			pemFromPriv := PrivateToPem(key)
			id, _ := shortid.Generate()
			_ = SaveToDirectory(path, id+".priv", pemFromPriv, globalPassphrase)
			slog.Info("Exported private Key", "path", filepath.Join(path, id+".priv"))
		} else if priv.Key != mainId.private {
			key := Key[*rsa.PrivateKey, crypto.PublicKey]{Public: pub, algType: RSA}
			pemFromPub := PublicToPem(key)
			id, _ := shortid.Generate()
			_ = SaveToDirectory(path, id+".pub", pemFromPub, nil)
			slog.Info("Exported Public Key", "path", filepath.Join(path, id+".pub"))
		}
	}
	return nil
}

func SaveToDirectory(path string, name string, pem []byte, passphrase *string) error {
	absolutePath := filepath.Clean(path)
	_, err := os.ReadDir(absolutePath)
	if err == nil {
		f, _ := os.Create(filepath.Join(absolutePath, name))
		if passphrase == nil {
			_, _ = f.Write(pem)
		} else {
			_, _ = f.Write(Encrypt(*passphrase, pem))
		}
		_ = f.Close()
	}
	return nil
}

func isTheMainIdentity(file os.DirEntry) bool {
	return strings.TrimSpace(file.Name()) == "main.priv"
}

func importMainIdentityPrivateKey(absolutePath string, file os.DirEntry, passphrase *string, wallet *Rsa) {
	mainId, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
	slog.Info("Reading main identity", "filepath", filepath.Join(absolutePath, file.Name()))
	if passphrase != nil {
		slog.Info("decrypting the main identity...")
		unencrypted := Decrypt(*passphrase, mainId)
		privKey, _ := PrivateFromPem(unencrypted)
		if isTheMainIdentity(file) {
			_ = wallet.SetMainIdentity(&privKey)
		}
		_ = wallet.Add(privKey)
	}
}

func importPublicKey(absolutePath string, file os.DirEntry, wallet *Rsa) {
	bytes, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
	fmt.Printf("Reading file %s\n", filepath.Join(absolutePath, file.Name()))
	key, _ := PublicFromPem(bytes)
	_ = wallet.Add(key)
}
