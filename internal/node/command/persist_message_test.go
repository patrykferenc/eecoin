package command_test

import (
	"testing"
	"time"

	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
	"github.com/stretchr/testify/assert"
)

func TestPersistMessageHandler_shouldNotCreate(t *testing.T) {
	handler, err := command.NewPersistMessageHandler(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, handler)
}

func TestPersistMessageHandler_shouldWork(t *testing.T) {
	// given
	repository := node.NewSimpleInFlightTransactionRepository()
	seen := node.NewSimpleSeenTransactionRepository()

	// and given
	handler, err := command.NewPersistMessageHandler(repository, seen)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	err = repository.Save(&node.Transaction{ID: "transaction-id", To: wallet.ID("to"), From: wallet.ID("from"), Timestamp: time.Now(), Content: "my silly message"})
	assert.NoError(t, err)

	// when
	err = handler.Handle(command.PersistMessage{TransactionID: "transaction-id"})

	// then
	assert.NoError(t, err)
	// and then
	wasSeen, err := seen.Seen("transaction-id")
	assert.NoError(t, err)
	assert.True(t, wasSeen)
	// and then
	transaction, err := repository.Get("transaction-id")
	assert.NoError(t, err)
	assert.Nil(t, transaction)
}
