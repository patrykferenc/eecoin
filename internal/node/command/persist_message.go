package command

import (
	"fmt"
	"log/slog"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

type PersistMessage struct {
	TransactionID blockchain.TransactionID
}

func NewPersistMessage(transactionID blockchain.TransactionID) (PersistMessage, error) {
	if transactionID == "" {
		return PersistMessage{}, fmt.Errorf("transactionID is required")
	}
	return PersistMessage{
		TransactionID: transactionID,
	}, nil
}

func (cmd PersistMessage) IsValid() error {
	if cmd == (PersistMessage{}) {
		return fmt.Errorf("command is required")
	}
	if cmd.TransactionID == "" {
		return fmt.Errorf("transactionID is required")
	}
	return nil
}

type PersistMessageHandler interface {
	Handle(cmd PersistMessage) error
}

type persistMessageHandler struct {
	repository node.InFlightTransactionRepository
	seen       node.SeenTransactionRepository
}

func (h *persistMessageHandler) Handle(cmd PersistMessage) error {
	if err := cmd.IsValid(); err != nil {
		return fmt.Errorf("can not persist message: invalid command: %w", err)
	}

	if err := h.seen.MarkSeen(cmd.TransactionID); err != nil {
		return fmt.Errorf("can not persist message: %w", err)
	}

	if err := h.repository.Discard(cmd.TransactionID); err != nil {
		return fmt.Errorf("can not persist message: %w", err)
	}

	slog.Info("Persisted message", "transactionID", cmd.TransactionID)

	return nil
}

func NewPersistMessageHandler(repository node.InFlightTransactionRepository, seen node.SeenTransactionRepository) (PersistMessageHandler, error) {
	if repository == nil {
		return nil, fmt.Errorf("repository is required")
	}
	if seen == nil {
		return nil, fmt.Errorf("seen is required")
	}
	return &persistMessageHandler{
		repository: repository,
		seen:       seen,
	}, nil
}
