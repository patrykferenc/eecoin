package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCoinbaseInvalid(t *testing.T) {
	assert := assert.New(t)
	// given a coinbase transaction
	in := NewInput(ID("some-output-id"), 0, "sign")
	tx := &Transaction{
		inputs:  []*Input{&in},
		outputs: []*Output{NewOutput(COINBASE_AMOUNT, "some-address")},
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

	// when validating the coinbase transaction
	err = validateCoinbase(tx, 1)

	// then no error should be returned
	assert.NoError(err)
}
