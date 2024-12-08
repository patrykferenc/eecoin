package wallet

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/teris-io/shortid"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type EcdsaKey = Key[*ecdsa.PrivateKey, crypto.PublicKey]
type EcdsaElem = KeyElement[ecdsa.PrivateKey]

type Ecdsa struct {
	MainId *EcdsaKey
	Keys   map[crypto.PublicKey]EcdsaElem
}

func (w *Ecdsa) SetMainIdentity(key *EcdsaKey) error {
	if key.private != nil {
		w.MainId = &EcdsaKey{private: key.private, Public: key.private.PublicKey, algType: ECDSA}
		return nil
	}
	return ErrPrivateKeyNotFound
}

func (w *Ecdsa) Add(key EcdsaKey) error {
	if key.private != nil {
		pub := key.private.Public()
		w.Keys[pub] = EcdsaElem{Key: key.private, Present: true}
		slog.Info("private-Public Key added to wallet", "pub", pub)
		return nil
	} else if key.Public != nil {
		w.Keys[key.Public] = EcdsaElem{Key: nil, Present: false}
		slog.Info("Public Key added to wallet", "pub", key.Public)
		return nil
	}
	return NoKeysFound
}

func (w *Ecdsa) Type() Algorithm {
	return ECDSA
}

func NewEcdsaKey() (EcdsaKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return EcdsaKey{algType: ECDSA}, nil
	}
	return EcdsaKey{private: key, Public: key.Public(), algType: ECDSA}, nil
}

func NewEcdsaWallet(mainId *EcdsaKey) *Ecdsa {
	wallet := &Ecdsa{Keys: make(map[crypto.PublicKey]EcdsaElem), MainId: mainId}
	_ = wallet.Add(*mainId)
	return wallet
}

func newEcdsaWalletWithoutId() *Ecdsa {
	wallet := &Ecdsa{Keys: make(map[crypto.PublicKey]EcdsaElem)}
	return wallet
}

func PublicFromPemEcdsa(pemData []byte) (EcdsaKey, error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return EcdsaKey{}, PemParseError
	}
	key, e := x509.ParsePKIXPublicKey(rest)
	if e != nil {
		slog.Error(e.Error())
		return EcdsaKey{}, PemParseError
	}
	return EcdsaKey{Public: key, algType: ECDSA}, nil
}

func PrivateFromPemEcdsa(pemData []byte) (EcdsaKey, error) {
	_, rest := pem.Decode(pemData)
	if rest == nil {
		return EcdsaKey{}, PemParseError
	}
	key, _ := x509.ParseECPrivateKey(rest)
	return EcdsaKey{private: key, Public: key.Public(), algType: ECDSA}, nil
}

func PrivateToPemEcdsa(key EcdsaKey) []byte {
	priv := key.private
	block, _ := x509.MarshalECPrivateKey(priv)
	return block
}

func PublicToPemEcdsa(key EcdsaKey) []byte {
	pub := key.Public
	block, _ := x509.MarshalPKIXPublicKey(pub)
	return block
}

func ReadWalletFromDirectoryEcdsa(path string, passphrase *string) (*Ecdsa, error) {

	absolutePath := filepath.Clean(path)
	directory, err := os.ReadDir(absolutePath)
	wallet := newEcdsaWalletWithoutId()

	if err != nil {
		return nil, err
	}

	for _, file := range directory {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pub") {
			slog.Info("Opened path", "path", file.Name())
			importPublicKeyEcdsa(absolutePath, file, wallet)
		} else if !file.IsDir() && strings.HasSuffix(file.Name(), ".priv") {
			slog.Info("Opened path", "path", file.Name())
			importPrivateKeyEcdsa(absolutePath, file, passphrase, wallet)
		}
	}
	return wallet, nil
}

func (w *Ecdsa) ExportWalletEcdsa(path string, globalPassphrase *string) error {
	mainId := *w.MainId
	mainIdPem := PrivateToPemEcdsa(mainId)
	_ = SaveToDirectoryEcdsa(path, "main.priv", mainIdPem, globalPassphrase)
	slog.Info("Exported main identity Key", "path", filepath.Join(path, "main.priv"))

	for pub, priv := range w.Keys {
		if priv.Present && priv.Key != mainId.private {
			key := EcdsaKey{private: priv.Key, algType: ECDSA}
			pemFromPriv := PrivateToPemEcdsa(key)
			id, _ := shortid.Generate()
			_ = SaveToDirectoryEcdsa(path, id+".priv", pemFromPriv, globalPassphrase)
			slog.Info("Exported private Key", "path", filepath.Join(path, id+".priv"))
		} else if priv.Key != mainId.private {
			key := EcdsaKey{Public: pub, algType: ECDSA}
			pemFromPub := PublicToPemEcdsa(key)
			id, _ := shortid.Generate()
			_ = SaveToDirectoryEcdsa(path, id+".pub", pemFromPub, nil)
			slog.Info("Exported Public Key", "path", filepath.Join(path, id+".pub"))
		}
	}
	return nil
}

func SaveToDirectoryEcdsa(path string, name string, pem []byte, passphrase *string) error {
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

func importPrivateKeyEcdsa(absolutePath string, file os.DirEntry, passphrase *string, wallet *Ecdsa) {
	mainId, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
	slog.Info("Reading main identity", "filepath", filepath.Join(absolutePath, file.Name()))
	if passphrase != nil {
		slog.Info("decrypting the main identity...")
		unencrypted := Decrypt(*passphrase, mainId)
		privKey, _ := PrivateFromPemEcdsa(unencrypted)
		if isTheMainIdentity(file) {
			_ = wallet.SetMainIdentity(&privKey)
		}
		_ = wallet.Add(privKey)
	}
}

func importPublicKeyEcdsa(absolutePath string, file os.DirEntry, wallet *Ecdsa) {
	bytes, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
	fmt.Printf("Reading file %s\n", filepath.Join(absolutePath, file.Name()))
	key, _ := PublicFromPemEcdsa(bytes)
	_ = wallet.Add(key)
}
