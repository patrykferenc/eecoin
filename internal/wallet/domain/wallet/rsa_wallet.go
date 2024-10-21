package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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
	return Key[*rsa.PrivateKey, crypto.PublicKey]{private: key, public: key.Public(), algType: RSA}, nil
}

func NewRsaWallet(mainId *Key[*rsa.PrivateKey, crypto.PublicKey]) *Rsa {
	wallet := &Rsa{keys: make(map[crypto.PublicKey]privateKeyElement), mainId: mainId}
	_ = wallet.Add(*mainId)
	return wallet
}

func newRsaWalletWithoutId() *Rsa {
	wallet := &Rsa{keys: make(map[crypto.PublicKey]privateKeyElement)}
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
	return Key[*rsa.PrivateKey, crypto.PublicKey]{public: key, algType: RSA}, nil
}

func PrivateFromPem(pemData []byte) (Key[*rsa.PrivateKey, crypto.PublicKey], error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return Key[*rsa.PrivateKey, crypto.PublicKey]{}, PemParseError
	}
	key, _ := x509.ParsePKCS1PrivateKey(rest)
	return Key[*rsa.PrivateKey, crypto.PublicKey]{private: key, public: key.Public(), algType: RSA}, nil
}

func PrivateToPem(key Key[*rsa.PrivateKey, crypto.PublicKey]) []byte {
	priv := key.private
	block := x509.MarshalPKCS1PrivateKey(priv)
	return block
}

func PublicToPem(key Key[*rsa.PrivateKey, crypto.PublicKey]) []byte {
	pub := key.public.(*rsa.PublicKey)
	block := x509.MarshalPKCS1PublicKey(pub)
	return block
}

func ReadWalletFromDirectory(path string, passphrase *string) (*Rsa, error) {

	absolutePath := filepath.Clean(path)
	directory, err := os.ReadDir(absolutePath)
	wallet := newRsaWalletWithoutId()

	if err == nil {
		for _, file := range directory {
			slog.Info("Opened path", "path", file.Name())
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".pub") {
				importPublicKey(absolutePath, file, wallet)
			} else if !file.IsDir() && strings.HasSuffix(file.Name(), ".priv") {
				importMainIdentityPrivateKey(absolutePath, file, passphrase, wallet)
			}
		}
		return wallet, nil
	}

	return nil, err
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
