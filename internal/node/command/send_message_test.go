package command

import (
	"testing"
	"time"

	"github.com/patrykferenc/eecoin/internal/common/event/eventtest"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/patrykferenc/eecoin/internal/wallet"
	"github.com/stretchr/testify/assert"
)

func TestSendMessageHandler_shouldNotCreate(t *testing.T) {
	handler, err := NewSendMessageHandler(nil, nil, nil, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, handler)
}

func TestSendMessageHandeler_shouldWork(t *testing.T) {
	// given
	repository := node.NewSimpleInFlightTransactionRepository()
	seen := node.NewSimpleSeenTransactionRepository()
	sender := &mockMessageSender{}
	peersToReturn := node.Peers{"localhost:2137", "localhost:1234"}
	peersRepo := &mockPeers{peers: peersToReturn, err: nil}
	publisher := eventtest.NewMockedPublisher()

	// and given
	handler, err := NewSendMessageHandler(repository, seen, sender, peersRepo, publisher)
	assert.NoError(t, err)
	assert.NotNil(t, handler)

	repository.Save(&node.Transaction{ID: "transaction-id", To: wallet.ID("to"), From: wallet.ID("from"), Timestamp: time.Now(), Content: "my silly message"})

	// when
	err = handler.Handle(SendMessage{TransactionID: "transaction-id"})

	// then
	assert.NoError(t, err)

	// and then
	assert.Len(t, publisher.Published(), 1)
}

type mockMessageSender struct {
	err error
}

func (m mockMessageSender) SendMessage(peers node.Peers, transaction *node.Transaction) error {
	return m.err
}

type mockPeers struct {
	peers node.Peers
	err   error
}

func (m mockPeers) Get() (node.Peers, error) {
	return m.peers, m.err
}
