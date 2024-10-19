package command_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/common/event/eventtest"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/stretchr/testify/assert"
)

func TestAcceptClientMessageHandler_shouldNotCreate(t *testing.T) {
	handler, err := command.NewAcceptClientMessageHandler(nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, handler)
}

func TestAcceptClientMessageHandler_shouldWork(t *testing.T) {
	publisher := eventtest.NewMockedPublisher()
	assert.Empty(t, publisher.Published())
	seen := node.NewSimpleSeenTransactionRepository()
	inFlight := node.NewSimpleInFlightTransactionRepository()
	handler, err := command.NewAcceptClientMessageHandler(inFlight, seen, publisher)
	assert.NoError(t, err)

	transactionID := blockchain.TransactionID(uuid.New().String())
	cmd, err := command.NewAcceptClientMessage(&node.Transaction{
		ID:      transactionID,
		Content: "content",
		From:    "from",
		To:      "to",
	})
	err = handler.Handle(cmd)
	assert.NoError(t, err)

	assert.NotEmpty(t, publisher.Published())
}

func TestAcceptClientMessageHandler_shouldNotPublish_whenTransactionAlreadySeen(t *testing.T) {
	// given
	publisher := eventtest.NewMockedPublisher()
	assert.Empty(t, publisher.Published())
	seen := node.NewSimpleSeenTransactionRepository()
	inFlight := node.NewSimpleInFlightTransactionRepository()
	handler, err := command.NewAcceptClientMessageHandler(inFlight, seen, publisher)
	assert.NoError(t, err)

	// and
	transactionID := blockchain.TransactionID(uuid.New().String())
	cmd, err := command.NewAcceptClientMessage(&node.Transaction{
		ID:      transactionID,
		Content: "content",
		From:    "from",
		To:      "to",
	})

	// when
	err = handler.Handle(cmd)

	// then
	assert.NoError(t, err)
	assert.Len(t, publisher.Published(), 1)

	// and given
	seen.MarkSeen(transactionID)
	// and when
	err = handler.Handle(cmd)
	// then
	assert.Error(t, err)
	assert.Len(t, publisher.Published(), 1)
}
