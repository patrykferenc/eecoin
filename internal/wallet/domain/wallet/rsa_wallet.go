package wallet

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
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

				handlePrivateKeyUnencrypted(absolutePath, file, wallet)

			} else if !file.IsDir() && strings.TrimSpace(file.Name()) == "main.priv" {

				mainId, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
				fmt.Printf("Reading main identity file %s\n", filepath.Join(absolutePath, file.Name()))
				// unencrypt
				if passphrase != nil {
					unencrypted := decrypt(*passphrase, mainId)
					mainKey, _ := PrivateFromPem(unencrypted)
					_ = wallet.SetMainIdentity(&mainKey)
				}
			}
		}
		return wallet, nil
	}

	return nil, err
}

func handlePrivateKeyUnencrypted(absolutePath string, file os.DirEntry, wallet *Rsa) {
	bytes, _ := os.ReadFile(filepath.Join(absolutePath, file.Name()))
	fmt.Printf("Reading file %s\n", filepath.Join(absolutePath, file.Name()))
	key, _ := PublicFromPem(bytes)
	_ = wallet.Add(key)
}

func SaveToDirectory(path string, name string, pem []byte, passphrase *string) error {

	absolutePath := filepath.Clean(path)
	_, err := os.ReadDir(absolutePath)
	if err == nil {
		f, _ := os.Create(filepath.Join(absolutePath, name))
		if passphrase == nil {
			_, _ = f.Write(pem)
		} else {
			_, _ = f.Write(encrypt(*passphrase, pem))
		}
		_ = f.Close()
	}

	return nil
}
func encrypt(passphrase string, plaintext []byte) []byte {
	key, salt := deriveKey(passphrase, nil)
	iv := make([]byte, 12)
	rand.Read(iv)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, plaintext, nil)
	return []byte(hex.EncodeToString(salt) + "kloszard" + hex.EncodeToString(iv) + "kloszard" + hex.EncodeToString(data))
}

func decrypt(passphrase string, ciphertext []byte) []byte {
	arr := strings.Split(string(ciphertext), "kloszard")
	salt, _ := hex.DecodeString(arr[0])
	iv, _ := hex.DecodeString(arr[1])
	data, _ := hex.DecodeString(arr[2])
	key, _ := deriveKey(passphrase, salt)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data, _ = aesgcm.Open(nil, iv, data, nil)
	return data
}
func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}
