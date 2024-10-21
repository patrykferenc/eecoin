package main

import (
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
)

func GenerateKeyToFile(path string, name string, passphrase string) error {
	key, _ := wallet.NewRsaKey()
	pubPem := wallet.PublicToPem(key)
	privPem := wallet.PrivateToPem(key)
	_ = wallet.SaveToDirectory(path, name+".priv", privPem, &passphrase)
	_ = wallet.SaveToDirectory(path, name+".pub", pubPem, nil)
	return nil
}
