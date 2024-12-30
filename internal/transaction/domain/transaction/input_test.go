package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignInput(t *testing.T) {
	assert := assert.New(t)

	// given
	outputID := ID("someOutputID")
	outputIndex := 1
	signatureOmmitedFromSigning := "irrelevant"
	input := NewInput(outputID, outputIndex, signatureOmmitedFromSigning)
	// and a signer
	privateSender, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	ourAddress, err := x509.MarshalPKIXPublicKey(privateSender.Public())
	require.NoError(t, err)
	// and an ID to sign
	idToSign := ID("someID") // TODO?
	// and an unspent output
	referencedOutput := UnspentOutput{
		outputID: outputID,
		address:  hex.EncodeToString(ourAddress),
	}

	// when
	err = input.sign(privateSender, idToSign, referencedOutput)
	assert.NoError(err)

	// then
	assert.NotEmpty(input.Signature())
	assert.NotEqual(input.Signature(), signatureOmmitedFromSigning)
}

func TestSignInputShouldNotSignWhenAddressDoesNotMatch(t *testing.T) {
	assert := assert.New(t)

	// given
	outputID := ID("someOutputID")
	outputIndex := 1
	signatureOmmitedFromSigning := "irrelevant"
	input := NewInput(outputID, outputIndex, signatureOmmitedFromSigning)
	// and a signer
	privateSender, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	// and an ID to sign
	idToSign := ID("someID")
	// and an unspent output
	referencedOutput := UnspentOutput{
		outputID: outputID,
		address:  "different address",
	}

	// when
	err = input.sign(privateSender, idToSign, referencedOutput)

	// then
	assert.Error(err)
	assert.Equal("output Addr does not match the signer Addr", err.Error())
}
