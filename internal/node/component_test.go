package node_test

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/common/mock"
	component "github.com/patrykferenc/eecoin/internal/node"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/stretchr/testify/assert"
)

// Scenario:
// Given a message from the client
// When the message is accepted by the node
// Then the message is broadcasted to the network
// Then the message is persisted on the blockchain
func TestMessaging_fromClient(t *testing.T) {
	assertThat := assert.New(t)

	// given a message from the client
	transaction := node.Transaction{
		ID:        blockchain.TransactionID(uuid.NewString()),
		Content:   "irrelevant content",
		Timestamp: time.Now(),
		From:      "client",
		To:        "client2",
	}

	// and given
	msg, err := command.NewAcceptClientMessage(&transaction)
	assertThat.NoError(err)

	// and given deps
	peers := mock.NewPeers(nil) // no peers needed
	seen := node.NewSimpleSeenTransactionRepository()
	repo := node.NewSimpleInFlightTransactionRepository()
	pub := event.NewChannelBroker()
	sender := mock.MessageSender{}
	component, err := component.NewComponent(pub, peers, seen, repo, sender)
	assertThat.NoError(err)
	acceptClientMessageHandler := component.Commands.AcceptClientMessage
	var wg sync.WaitGroup
	wg.Add(2)

	pub.Route("x.message.send", func(e event.Event) error {
		defer wg.Done()
		data, ok := e.Data().(node.SendMessageEvent)
		if !ok {
			assertThat.Fail("invalid event data")
		}
		assertThat.Equal(transaction.ID, data.TransactionID)
		// and then
		trans, err := repo.Get(transaction.ID)
		assertThat.NoError(err)
		assertThat.NotNil(trans)

		// and when
		cmd, err := command.NewSendMessage(data.TransactionID)
		assertThat.NoError(err)
		err = component.Commands.SendMessage.Handle(cmd)
		assertThat.NoError(err)

		return nil
	})
	pub.Route("x.message.sent", func(e event.Event) error {
		defer wg.Done()
		data, ok := e.Data().(node.MessageSentEvent)
		if !ok {
			assertThat.Fail("invalid event data")
		}
		assertThat.Equal(transaction.ID, data.TransactionID)
		// and given
		cmd, err := command.NewPersistMessage(data.TransactionID)
		assertThat.NoError(err)

		// when
		err = component.Commands.PersistMessage.Handle(cmd)
		assertThat.NoError(err)

		// then
		actualInFlight, err := repo.Get(transaction.ID)
		assertThat.NoError(err)
		assertThat.Nil(actualInFlight) // deleted from in flight
		// and then
		actualSeen, err := seen.Seen(transaction.ID)
		assertThat.NoError(err)
		assertThat.True(actualSeen) // added to seen
		return nil
	})

	// when
	err = acceptClientMessageHandler.Handle(msg)

	// then
	assertThat.NoError(err)

	// and then
	wg.Wait()
}
