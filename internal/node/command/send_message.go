package command

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

type SendMessage struct {
	TransactionID blockchain.TransactionID
}

func NewSendMessage(transactionID blockchain.TransactionID) (SendMessage, error) {
	if transactionID == "" {
		return SendMessage{}, fmt.Errorf("transactionID is required")
	}
	return SendMessage{
		TransactionID: transactionID,
	}, nil
}

func (cmd SendMessage) IsValid() error {
	if cmd == (SendMessage{}) {
		return fmt.Errorf("command is required")
	}
	if cmd.TransactionID == "" {
		return fmt.Errorf("transactionID is required")
	}
	return nil
}

type SendMessageHandler interface {
	Handle(cmd SendMessage) error
}

type sendMessageHandler struct {
	repository       node.InFlightTransactionRepository
	seen             node.SeenTransactionRepository
	sender           node.MessageSender
	peeersRepository node.PeersRepository
	publisher        event.Publisher
}

func (h *sendMessageHandler) Handle(cmd SendMessage) error {
	if err := cmd.IsValid(); err != nil {
		return fmt.Errorf("can not send message: invalid command: %w", err)
	}

	if seen, err := h.seen.Seen(cmd.TransactionID); err != nil {
		return fmt.Errorf("can not send message: %w", err)
	} else if seen {
		return fmt.Errorf("can not send message: transaction already seen")
	}

	transaction, err := h.repository.Get(cmd.TransactionID)
	if err != nil {
		return fmt.Errorf("can not send message: %w", err)
	}

	peers, err := h.peeersRepository.Get()
	if err != nil {
		return fmt.Errorf("can not send message: %w", err)
	}

	err = h.sender.SendMessage(peers, transaction)
	if err != nil {
		return fmt.Errorf("error when sending message: %w", err)
	}

	event, err := event.New(&node.MessageSentEvent{TransactionID: cmd.TransactionID})
	if err != nil {
		return fmt.Errorf("can not send message: %w", err)
	}

	err = h.publisher.Publish(event)
	if err != nil {
		return fmt.Errorf("can not send message: %w", err)
	}

	return nil
}

func NewSendMessageHandler(
	repository node.InFlightTransactionRepository,
	seen node.SeenTransactionRepository,
	sender node.MessageSender,
	peersRepository node.PeersRepository,
	publisher event.Publisher,
) (SendMessageHandler, error) {
	if repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	if seen == nil {
		return nil, fmt.Errorf("seen is required")
	}
	if sender == nil {
		return nil, fmt.Errorf("sender is required")
	}
	if peersRepository == nil {
		return nil, fmt.Errorf("peersRepository is required")
	}
	if publisher == nil {
		return nil, fmt.Errorf("publisher is required")
	}

	return &sendMessageHandler{
		repository:       repository,
		seen:             seen,
		sender:           sender,
		peeersRepository: peersRepository,
		publisher:        publisher,
	}, nil
}
