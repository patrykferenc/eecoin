package wallet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	mainId, _ := NewRsaKey()
	resultToPem := PrivateToPemRsa(mainId)

	encrypted := Encrypt("dupa", resultToPem)
	decrypted := Decrypt("dupa", encrypted)
	assert.Equal(t, resultToPem, decrypted)
}
