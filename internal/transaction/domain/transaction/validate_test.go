package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCoinbaseInvalid(t *testing.T) {
	assert := assert.New(t)
	// given a coinbase transaction
	in := NewInput(ID("some-output-id"), 0, "sign")
	tx := &Transaction{
		In: []*Input{&in},
		Ou: []*Output{NewOutput(COINBASE_AMOUNT, "some-address")},
	}
	blockHeight := 1

	// when validating the coinbase transaction
	err := validateCoinbase(tx, blockHeight)

	// then no error should be returned
	assert.Error(err)
}

func TestValidateCoinbaseValid(t *testing.T) {
	assert := assert.New(t)
	// given a coinbase transaction
	tx, err := NewCoinbase("some-address", 1)
	assert.NoError(err)

	// when validating the coinbase transaction
	err = validateCoinbase(tx, 1)

	// then no error should be returned
	assert.NoError(err)
}

func TestValidateTransactionIn(t *testing.T) {
	assert := assert.New(t)

	// --- Setup for both cases ---

	// Generate key pair for the sender
	privateSender, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(err)
	senderAddrRaw, err := x509.MarshalPKIXPublicKey(privateSender.Public())
	assert.NoError(err)
	senderAddr := hex.EncodeToString(senderAddrRaw)

	// Create a valid unspent output
	unspentOutput := NewUnspentOutput("some-tx-id", 0, 100, senderAddr)

	// Mock UnspentOutputRepository
	mockUnspentRepo := &mockUnspentOutputRepository{
		// UnspentOutputs: map[string][]UnspentOutput{
		// 	string(senderAddr): {unspentOutput},
		// },
		UnspentOutputs: map[string][]UnspentOutput{
			string(senderAddr): {unspentOutput},
		},
	}

	// Create a transaction referencing the unspent output
	tx := &Transaction{
		Id: "some-tx-id",
	}

	// Valid signature
	txIDHash := sha256.Sum256([]byte(tx.ID().String()))
	r, s, err := ecdsa.Sign(rand.Reader, privateSender, txIDHash[:])
	assert.NoError(err)

	signature := append(r.Bytes(), s.Bytes()...)

	validInput := &Input{
		OutputId:  "some-tx-id",
		OutputIdx: 0,
		Sig:       hex.EncodeToString(signature),
	}

	// --- Valid Case ---

	t.Run("ValidTransactionInput", func(t *testing.T) {
		err := validateTransactionIn(validInput, tx, mockUnspentRepo)
		assert.NoError(err, "expected no error for valid input")
	})

	// --- Invalid Case (tampered signature) ---

	t.Run("InvalidTransactionInput", func(t *testing.T) {
		// Create a tampered signature (by modifying one byte)
		tamperedSignature := append(r.Bytes(), s.Bytes()...)
		tamperedSignature[len(tamperedSignature)-1] ^= 0xFF // Flip the last byte

		invalidInput := &Input{
			OutputId:  "some-tx-id",
			OutputIdx: 0,
			Sig:       hex.EncodeToString(tamperedSignature),
		}

		err := validateTransactionIn(invalidInput, tx, mockUnspentRepo)
		assert.Error(err, "expected error for invalid input")
		assert.Contains(err.Error(), "invalid signature", "error should indicate invalid signature")
	})
}

type mockUnspentOutputRepository struct {
	Called         int
	UnspentOutputs map[string][]UnspentOutput
}

func (r *mockUnspentOutputRepository) GetByAddress(address string) ([]UnspentOutput, error) {
	r.Called++
	return r.UnspentOutputs[address], nil
}

func (r *mockUnspentOutputRepository) GetAll() ([]UnspentOutput, error) {
	r.Called++
	var uos []UnspentOutput
	for _, outputs := range r.UnspentOutputs {
		uos = append(uos, outputs...)
	}
	return uos, nil
}

func (r *mockUnspentOutputRepository) Set(_ []UnspentOutput) error {
	panic("not implemented")
}

func (r *mockUnspentOutputRepository) GetByOutputIDAndIndex(outputID ID, outputIndex int) (UnspentOutput, error) {
	r.Called++
	for _, outputs := range r.UnspentOutputs {
		for _, output := range outputs {
			if output.OutputID() == outputID && output.OutputIndex() == outputIndex {
				return output, nil
			}
		}
	}
	return UnspentOutput{}, fmt.Errorf("output with ID %s and index %d not found", outputID, outputIndex)
}
