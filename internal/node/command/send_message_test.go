package command_test

import (
	"testing"
	"time"

	"github.com/patrykferenc/eecoin/internal/common/event/eventtest"
	"github.com/patrykferenc/eecoin/internal/common/mock"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
	"github.com/stretchr/testify/assert"
)

func TestSendMessageHandler_shouldNotCreate(t *testing.T) {
	handler, err := command.NewSendMessageHandler(nil, nil, nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, handler)
}

func TestSendMessageHandeler_shouldWork(t *testing.T) {
	// given
	repository := node.NewSimpleInFlightTransactionRepository()
	seen := node.NewSimpleSeenTransactionRepository()
	sender := mock.MessageSender{}
	peersToReturn := []string{"localhost:2137", "localhost:1234"}
	peersRepo := mock.Peers{Peers: peersToReturn}
	publisher := eventtest.NewMockedPublisher()

	// and given
	handler, err := command.NewSendMessageHandler(repository, seen, sender, peersRepo, publisher)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	err = repository.Save(&node.Transaction{ID: "transaction-id", To: wallet.ID("to"), From: wallet.ID("from"), Timestamp: time.Now(), Content: "my silly message"})
	assert.NoError(t, err)

	// when
	err = handler.Handle(command.SendMessage{TransactionID: "transaction-id"})

	// then
	assert.NoError(t, err)

	// and then
	assert.Len(t, publisher.Published(), 1)
}
