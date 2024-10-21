package wallet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRsaKey(t *testing.T) {
	assertThat := assert.New(t)

	//given
	rsaKeySizeInBits := 2048
	bitsInByte := 8
	rsaKey, _ := rsa.GenerateKey(rand.Reader, rsaKeySizeInBits)
	pubFromKey := rsaKey.Public()

	//when
	result, err := NewRsaKey()

	//then
	assertThat.NotNil(result)
	assertThat.Nil(err)
	assertThat.IsTypef(result.private, rsaKey, "Private Keys should be the same type %T", rsaKey)
	assertThat.IsTypef(result.Public, pubFromKey, "Private Keys should be the same type %T", pubFromKey)
	assertThat.Equal(result.private.Size(), rsaKeySizeInBits/bitsInByte)
}

func TestNewRsaWallet(t *testing.T) {
	assertThat := assert.New(t)

	//given
	mainId, _ := NewRsaKey()
	//when
	result := NewRsaWallet(&mainId)
	//then
	assertThat.NotNil(result)
	assertThat.Equal(&mainId, result.MainId)
	assertThat.NotNil(result.Keys)
}

func TestRsaWallet_Add(t *testing.T) {
	assertThat := assert.New(t)

	//given
	mainId, _ := NewRsaKey()
	key1, _ := NewRsaKey()
	key2, _ := NewRsaKey()

	justPrivate := Key[*rsa.PrivateKey, crypto.PublicKey]{private: key1.private, algType: RSA}
	justPublic := Key[*rsa.PrivateKey, crypto.PublicKey]{Public: key2.Public, algType: RSA}
	allNil := Key[*rsa.PrivateKey, crypto.PublicKey]{algType: RSA}
	//when - then
	result := NewRsaWallet(&mainId)
	_ = result.Add(justPrivate)
	assertThat.Len(result.Keys, 2)
	assertThat.True(result.Keys[key1.Public].Present)
	//when - then
	_ = result.Add(justPublic)
	assertThat.Len(result.Keys, 3)
	assertThat.False(result.Keys[justPublic.Public].Present)
	//when - then
	err := result.Add(allNil)
	assertThat.Len(result.Keys, 3)
	assertThat.ErrorIs(err, NoKeysFound)

}
func TestRsaWallet_Type(t *testing.T) {
	assertThat := assert.New(t)

	//given
	mainId, _ := NewRsaKey()
	//when
	result := NewRsaWallet(&mainId)
	//then
	assertThat.Equal(RSA, result.Type())
}
func TestRsaWallet_SetMainIdentity(t *testing.T) {
	assertThat := assert.New(t)

	//given
	mainId, _ := NewRsaKey()
	anotherId, _ := NewRsaKey()
	//when
	result := NewRsaWallet(&mainId)
	//then
	_ = result.SetMainIdentity(&anotherId)
	assertThat.Equal(anotherId.private, result.MainId.private)

}

func TestToPem(t *testing.T) {
	assertThat := assert.New(t)
	mainId, _ := NewRsaKey()

	resultToPem := PublicToPem(mainId)
	resultFromPem, _ := PublicFromPem(resultToPem)

	resultPrivToPem := PrivateToPem(mainId)
	resultPrivFromPem, _ := PrivateFromPem(resultPrivToPem)

	assertThat.Equal(resultFromPem.Public, mainId.Public)
	assertThat.Equal(resultPrivFromPem.private, mainId.private)
}
func TestSavingMainIdentity(t *testing.T) {
	assertThat := assert.New(t)
	//given
	mainId, _ := NewRsaKey()
	privOnly, _ := NewRsaKey()
	pubOnly, _ := NewRsaKey()

	pass := "dupa"
	w := NewRsaWallet(&mainId)
	w.Add(privOnly)
	w.Add(pubOnly)
	//when
	saveError := w.ExportWallet("/tmp", &pass)

	wallet, readError := ReadWalletFromDirectory("/tmp", &pass)
	//then
	assertThat.Nil(saveError)
	assertThat.Nil(readError)
	assertThat.NotNil(wallet)
	assertThat.NotNil(wallet.Keys)
	assertThat.NotNil(wallet.MainId)
	assertThat.Equal(mainId.private, wallet.MainId.private)
	assertThat.Len(wallet.Keys, 3)
}
