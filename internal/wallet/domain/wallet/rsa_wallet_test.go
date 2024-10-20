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
	assertThat.IsTypef(result.private, rsaKey, "Private keys should be the same type %T", rsaKey)
	assertThat.IsTypef(result.public, pubFromKey, "Private keys should be the same type %T", pubFromKey)
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
	assertThat.Equal(&mainId, result.mainId)
	assertThat.NotNil(result.keys)
}

func TestRsaWallet_Add(t *testing.T) {
	assertThat := assert.New(t)

	//given
	mainId, _ := NewRsaKey()
	key1, _ := NewRsaKey()
	key2, _ := NewRsaKey()

	justPrivate := Key[*rsa.PrivateKey, crypto.PublicKey]{private: key1.private, algType: RSA}
	justPublic := Key[*rsa.PrivateKey, crypto.PublicKey]{public: key2.public, algType: RSA}
	allNil := Key[*rsa.PrivateKey, crypto.PublicKey]{algType: RSA}
	//when - then
	result := NewRsaWallet(&mainId)
	_ = result.Add(justPrivate)
	assertThat.Len(result.keys, 2)
	assertThat.True(result.keys[key1.public].present)
	//when - then
	_ = result.Add(justPublic)
	assertThat.Len(result.keys, 3)
	assertThat.False(result.keys[justPublic.public].present)
	//when - then
	err := result.Add(allNil)
	assertThat.Len(result.keys, 3)
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
	assertThat.Equal(anotherId.private, result.mainId.private)

}

func TestToPem(t *testing.T) {
	assertThat := assert.New(t)
	mainId, _ := NewRsaKey()

	resultToPem := ToPem(mainId)
	resultFromPem, _ := FromPem(resultToPem)

	assertThat.Equal(resultFromPem.public, mainId.public)
}
