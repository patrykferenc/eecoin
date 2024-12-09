package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateEcdsaKey(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	//given
	ecdsaKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubFromKey := ecdsaKey.Public()

	//when
	result, err := NewEcdsaKey()

	//then
	assertThat.NotNil(result)
	assertThat.Nil(err)
	assertThat.IsTypef(result.private, ecdsaKey, "Private Keys should be the same type %T", ecdsaKey)
	assertThat.IsTypef(result.Public, pubFromKey, "Private Keys should be the same type %T", pubFromKey)
}

func TestNewEcdsaWallet(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	//given
	mainId, _ := NewEcdsaKey()
	//when
	result := NewEcdsaWallet(&mainId)
	//then
	assertThat.NotNil(result)
	assertThat.Equal(&mainId, result.MainId)
	assertThat.NotNil(result.Keys)
}

func TestEcdsaWallet_Add(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	//given
	mainId, _ := NewEcdsaKey()
	key1, _ := NewEcdsaKey()
	key2, _ := NewEcdsaKey()

	justPrivate := EcdsaKey{private: key1.private, algType: RSA}
	justPublic := EcdsaKey{Public: key2.Public, algType: RSA}
	allNil := EcdsaKey{algType: RSA}
	//when - then
	result := NewEcdsaWallet(&mainId)
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
func TestEcdsaWallet_Type(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	//given
	mainId, _ := NewEcdsaKey()
	//when
	result := NewEcdsaWallet(&mainId)
	//then
	assertThat.Equal(ECDSA, result.Type())
}
func TestEcdsaWallet_SetMainIdentity(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	//given
	mainId, _ := NewEcdsaKey()
	anotherId, _ := NewEcdsaKey()
	//when
	result := NewEcdsaWallet(&mainId)
	//then
	_ = result.SetMainIdentity(&anotherId)
	assertThat.Equal(anotherId.private, result.MainId.private)

}
func TestToPemEcdsa(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)
	mainId, _ := NewEcdsaKey()

	resultToPem := PublicToPemEcdsa(mainId)
	resultFromPem, _ := PublicFromPemEcdsa(resultToPem)

	resultPrivToPem := PrivateToPemEcdsa(mainId)
	resultPrivFromPem, _ := PrivateFromPemEcdsa(resultPrivToPem)

	assertThat.Equal(resultFromPem.Public, mainId.Public)
	assertThat.Equal(resultPrivFromPem.private, mainId.private)
}
func TestSavingMainIdentityEcdsa(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)
	dir := t.TempDir()
	//given
	mainId, _ := NewEcdsaKey()
	privOnly, _ := NewEcdsaKey()
	pubOnly, _ := NewEcdsaKey()

	pass := "dupa"
	w := NewEcdsaWallet(&mainId)
	_ = w.Add(privOnly)
	_ = w.Add(pubOnly)
	//when
	saveError := w.ExportWalletEcdsa(dir, &pass)

	wallet, readError := ReadWalletFromDirectoryEcdsa(dir, &pass)
	//then
	assertThat.Nil(saveError)
	assertThat.Nil(readError)
	assertThat.NotNil(wallet)
	assertThat.NotNil(wallet.Keys)
	assertThat.NotNil(wallet.MainId)
	assertThat.Equal(mainId.private, wallet.MainId.private)
	assertThat.Len(wallet.Keys, 3)
}
