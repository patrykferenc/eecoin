package command

import (
	"errors"
	"fmt"

	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

type AcceptClientMessage struct {
	transaction *node.Transaction // TODO: can potentially be refactored to just plain fields?
}

func NewAcceptClientMessage(transaction *node.Transaction) (AcceptClientMessage, error) {
	if transaction == nil {
		return AcceptClientMessage{}, errors.New("transaction is required")
	}
	return AcceptClientMessage{
		transaction: transaction,
	}, nil
}

func (cmd AcceptClientMessage) IsValid() error {
	if cmd.transaction == nil {
		return errors.New("transaction is required")
	}
	return nil
}

func (cmd AcceptClientMessage) Transaction() *node.Transaction {
	return cmd.transaction
}

type AcceptClientMessageHandler interface {
	Handle(cmd AcceptClientMessage) error
}

type acceptClientMessageHandler struct {
	inFlightRepository node.InFlightTransactionRepository
	seenRepository     node.SeenTransactionRepository
	publisher          event.Publisher
}

func (h *acceptClientMessageHandler) Handle(cmd AcceptClientMessage) error {
	if err := cmd.IsValid(); err != nil {
		return fmt.Errorf("can not accept message: invalid command: %w", err)
	}

	if seen, err := h.seenRepository.Seen(cmd.Transaction().ID); err != nil {
		return fmt.Errorf("can not accept message: seen check failed: %w", err)
	} else if seen {
		return fmt.Errorf("can not accept message: transaction already seen")
	}

	if err := h.inFlightRepository.Save(cmd.Transaction()); err != nil {
		return fmt.Errorf("can not accept message: failed to save transaction: %w", err)
	}

	event, err := event.New(node.SendMessageEvent{TransactionID: cmd.Transaction().ID}, "x.message.send")
	if err != nil {
		return fmt.Errorf("can not accept message: failed to create event: %w", err)
	}
	if err := h.publisher.Publish(event); err != nil {
		return fmt.Errorf("can not accept message: failed to publish event: %w", err)
	}

	return nil
}

func NewAcceptClientMessageHandler(
	inFlightRepository node.InFlightTransactionRepository,
	seenRepository node.SeenTransactionRepository,
	publisher event.Publisher,
) (AcceptClientMessageHandler, error) {
	if inFlightRepository == nil {
		return nil, errors.New("inFlightRepository is required")
	}
	if seenRepository == nil {
		return nil, errors.New("seenRepository is required")
	}
	if publisher == nil {
		return nil, errors.New("publisher is required")
	}

	return &acceptClientMessageHandler{
		inFlightRepository: inFlightRepository,
		seenRepository:     seenRepository,
		publisher:          publisher,
	}, nil
}
