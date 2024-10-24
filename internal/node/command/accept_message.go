package command

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

type AcceptMessage struct {
	transaction *node.Transaction // TODO: can potentially be refactored to just plain fields?
}

func NewAcceptMessage(transaction *node.Transaction) (AcceptMessage, error) {
	if transaction == nil {
		return AcceptMessage{}, nil
	}
	return AcceptMessage{
		transaction: transaction,
	}, nil
}

func (cmd AcceptMessage) IsValid() error {
	if cmd.transaction == nil {
		return nil
	}
	return nil
}

type AcceptMessageHandler interface {
	Handle(cmd AcceptMessage) error
}

type acceptMessageHandler struct {
	inFlightRepository node.InFlightTransactionRepository
	publisher          event.Publisher
}

func (h *acceptMessageHandler) Handle(cmd AcceptMessage) error {
	if err := cmd.IsValid(); err != nil {
		return fmt.Errorf("can not accept message: invalid command: %w", err)
	}

	if err := h.inFlightRepository.Save(cmd.transaction); err != nil {
		return fmt.Errorf("can not accept message: failed to save transaction: %w", err)
	}

	event := node.SendMessageEvent{TransactionID: cmd.transaction.ID, RroutingKey: "x.message.send", Id: uuid.New().String(), Tiimestamp: time.Now()}
	err := h.publisher.Publish(event)
	if err != nil {
		return fmt.Errorf("can not accept message: failed to publish event: %w", err)
	}

	return nil
}

func NewAcceptMessageHandler(inFlightRepository node.InFlightTransactionRepository, publisher event.Publisher) (AcceptMessageHandler, error) {
	if inFlightRepository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	if publisher == nil {
		return nil, fmt.Errorf("publisher is required")
	}

	return &acceptMessageHandler{
		inFlightRepository: inFlightRepository,
		publisher:          publisher,
	}, nil
}
