package command_test

import (
	"crypto/sha256"
	"testing"

	"github.com/patrykferenc/eecoin/internal/common/mock"
	"github.com/patrykferenc/eecoin/internal/transaction/command"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldAddTransaction(t *testing.T) {
	assert := assert.New(t)
	// given
	poolRepository := mock.NewPoolRepository()
	pool := transaction.NewPool(poolRepository)
	publisher := &mock.Publisher{}
	handler := command.NewAddTransactionHandler(publisher, pool)

	// and given inputs
	someInput := transaction.NewInput(transaction.ID("output-id"), 1, "signature")
	inputs := []*transaction.Input{
		&someInput,
	}
	someOutput := transaction.NewOutput(10, "addressTo")
	outputs := []*transaction.Output{
		someOutput,
	}
	// and given matching hash
	h := sha256.New()
	_, err := h.Write([]byte("output-id110addressTo"))
	require.NoError(t, err)

	// when
	err = handler.Handle(command.AddTransaction{
		ProvidedID: string(h.Sum(nil)),
		Inputs:     inputs,
		Outputs:    outputs,
	})

	// then
	assert.NoError(err)
	assert.Equal(1, poolRepository.Called)
	assert.Equal(1, publisher.Called)
}

func TestShouldNotAddInvalidTransaction(t *testing.T) {
	assert := assert.New(t)
	// given
	poolRepository := mock.NewPoolRepository()
	pool := transaction.NewPool(poolRepository)
	publisher := &mock.Publisher{}
	handler := command.NewAddTransactionHandler(publisher, pool)

	// and given inputs
	someInput := transaction.NewInput(transaction.ID("output-id"), 1, "signature")
	inputs := []*transaction.Input{
		&someInput,
	}
	someOutput := transaction.NewOutput(10, "addressTo")
	outputs := []*transaction.Output{
		someOutput,
	}
	// and given non-matching hash
	h := sha256.New()
	_, err := h.Write([]byte("output-id110addressTo"))
	require.NoError(t, err)

	// when
	err = handler.Handle(command.AddTransaction{
		ProvidedID: "invalid",
		Inputs:     inputs,
		Outputs:    outputs,
	})

	// then
	assert.Error(err)
	assert.Equal(0, poolRepository.Called)
	assert.Equal(0, publisher.Called)
}
