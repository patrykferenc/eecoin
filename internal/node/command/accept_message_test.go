package command_test

import (
	"testing"
	"time"

	"github.com/patrykferenc/eecoin/internal/common/event/eventtest"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/stretchr/testify/assert"
)

func TestAcceptMessageHandler_shouldNotCreate(t *testing.T) {
	handler, err := command.NewAcceptMessageHandler(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, handler)
}

func TestAcceptMessageHandler_shouldWork(t *testing.T) {
	// given
	repository := node.NewSimpleInFlightTransactionRepository()
	publisher := eventtest.NewMockedPublisher()

	// and given
	handler, err := command.NewAcceptMessageHandler(repository, publisher)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	// and given
	transaction := &node.Transaction{
		ID:        "transaction-id",
		Content:   "my silliest message ever",
		Timestamp: time.Now(),
		From:      "dziabaducha",
		To:        "szniobert-okrutnik",
	}
	cmd, err := command.NewAcceptMessage(transaction)
	assert.NoError(t, err)
	// when
	err = handler.Handle(cmd)

	// then
	assert.NoError(t, err)
	// and then
	actual, err := repository.Get("transaction-id")
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	// and then
	assert.Len(t, publisher.Published(), 1)
}
