package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

func Encrypt(passphrase string, plaintext []byte) []byte {
	key, salt := deriveKey(passphrase, nil)
	iv := make([]byte, 12)
	rand.Read(iv)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, plaintext, nil)
	return []byte(hex.EncodeToString(salt) + "kloszard" + hex.EncodeToString(iv) + "kloszard" + hex.EncodeToString(data))
}

func Decrypt(passphrase string, ciphertext []byte) []byte {
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
