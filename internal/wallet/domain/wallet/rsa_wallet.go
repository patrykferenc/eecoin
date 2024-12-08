package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/teris-io/shortid"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const defaultRsaBitSize = 2048

type RsaKey = Key[*rsa.PrivateKey, crypto.PublicKey]
type RsaElem = KeyElement[rsa.PrivateKey]

type Rsa struct {
	MainId *RsaKey
	Keys   map[crypto.PublicKey]RsaElem
}

func (w *Rsa) SetMainIdentity(key *RsaKey) error {
	if key.private != nil {
		w.MainId = &RsaKey{private: key.private, Public: key.private.PublicKey, algType: RSA}
		return nil
	}
	return ErrPrivateKeyNotFound
}

func (w *Rsa) Add(key RsaKey) error {
	if key.private != nil {
		pub := key.private.Public()
		w.Keys[pub] = RsaElem{Key: key.private, Present: true}
		slog.Info("private-Public Key added to wallet", "pub", pub)
		return nil
	} else if key.Public != nil {
		w.Keys[key.Public] = RsaElem{Key: nil, Present: false}
		slog.Info("Public Key added to wallet", "pub", key.Public)
		return nil
	}
	return NoKeysFound
}
func (w *Rsa) Type() Algorithm {
	return RSA
}

func NewRsaKey() (RsaKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, defaultRsaBitSize)
	if err != nil {
		return RsaKey{algType: RSA}, nil
	}
	return RsaKey{private: key, Public: key.Public(), algType: RSA}, nil
}

func NewRsaWallet(mainId *RsaKey) *Rsa {
	wallet := &Rsa{Keys: make(map[crypto.PublicKey]RsaElem), MainId: mainId}
	_ = wallet.Add(*mainId)
	return wallet
}

func newRsaWalletWithoutId() *Rsa {
	wallet := &Rsa{Keys: make(map[crypto.PublicKey]RsaElem)}
	return wallet
}

func PublicFromPemRsa(pemData []byte) (RsaKey, error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return RsaKey{}, PemParseError
	}
	key, e := x509.ParsePKCS1PublicKey(rest)
	if e != nil {
		slog.Error(e.Error())
		return RsaKey{}, PemParseError
	}
	return RsaKey{Public: key, algType: RSA}, nil
}

func PrivateFromPemRsa(pemData []byte) (RsaKey, error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return RsaKey{}, PemParseError
	}
	key, _ := x509.ParsePKCS1PrivateKey(rest)
	return RsaKey{private: key, Public: key.Public(), algType: RSA}, nil
}

func PrivateToPemRsa(key RsaKey) []byte {
	priv := key.private
	block := x509.MarshalPKCS1PrivateKey(priv)
	return block
}

func PublicToPemRsa(key RsaKey) []byte {
	pub := key.Public.(*rsa.PublicKey)
	block := x509.MarshalPKCS1PublicKey(pub)
	return block
}

func ReadWalletFromDirectoryRsa(path string, passphrase *string) (*Rsa, error) {

	absolutePath := filepath.Clean(path)
	directory, err := os.ReadDir(absolutePath)
	wallet := newRsaWalletWithoutId()

	if err != nil {
		return nil, err
	}

	for _, file := range directory {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pub") {
			slog.Info("Opened path", "path", file.Name())
			importPublicKeyRsa(absolutePath, file, wallet)
		} else if !file.IsDir() && strings.HasSuffix(file.Name(), ".priv") {
			slog.Info("Opened path", "path", file.Name())
			importPrivateKeyRsa(absolutePath, file, passphrase, wallet)
		}
	}
	return wallet, nil
}

func (w *Rsa) ExportWalletRsa(path string, globalPassphrase *string) error {
	mainId := *w.MainId
	mainIdPem := PrivateToPemRsa(mainId)
	_ = SaveToDirectoryRsa(path, "main.priv", mainIdPem, globalPassphrase)
	slog.Info("Exported main identity Key", "path", filepath.Join(path, "main.priv"))

	for pub, priv := range w.Keys {
		if priv.Present && priv.Key != mainId.private {
			key := RsaKey{private: priv.Key, algType: RSA}
			pemFromPriv := PrivateToPemRsa(key)
			id, _ := shortid.Generate()
			_ = SaveToDirectoryRsa(path, id+".priv", pemFromPriv, globalPassphrase)
			slog.Info("Exported private Key", "path", filepath.Join(path, id+".priv"))
		} else if priv.Key != mainId.private {
			key := RsaKey{Public: pub, algType: RSA}
			pemFromPub := PublicToPemRsa(key)
			id, _ := shortid.Generate()
			_ = SaveToDirectoryRsa(path, id+".pub", pemFromPub, nil)
			slog.Info("Exported Public Key", "path", filepath.Join(path, id+".pub"))
		}
	}
	return nil
}

func SaveToDirectoryRsa(path string, name string, pem []byte, passphrase *string) error {
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

func importPrivateKeyRsa(absolutePath string, file os.DirEntry, passphrase *string, wallet *Rsa) {
	mainId, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
	slog.Info("Reading main identity", "filepath", filepath.Join(absolutePath, file.Name()))
	if passphrase != nil {
		slog.Info("decrypting the main identity...")
		unencrypted := Decrypt(*passphrase, mainId)
		privKey, _ := PrivateFromPemRsa(unencrypted)
		if isTheMainIdentity(file) {
			_ = wallet.SetMainIdentity(&privKey)
		}
		_ = wallet.Add(privKey)
	}
}

func importPublicKeyRsa(absolutePath string, file os.DirEntry, wallet *Rsa) {
	bytes, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
	fmt.Printf("Reading file %s\n", filepath.Join(absolutePath, file.Name()))
	key, _ := PublicFromPemRsa(bytes)
	_ = wallet.Add(key)
}
