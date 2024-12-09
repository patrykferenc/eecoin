package command

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type BroadcastTransaction struct {
	TransactionID string
}

type BroadcastTransactionHandler interface {
	Handle(BroadcastTransaction) error
}

type Broadcaster interface {
	Broadcast(transaction.Transaction) error
}

type broadcastTransactionHandler struct {
	pool        *transaction.Pool
	broadcaster Broadcaster
}

func NewBroadcastTransactionHandler(
	pool *transaction.Pool,
	broadcaster Broadcaster,
) BroadcastTransactionHandler {
	return &broadcastTransactionHandler{
		pool:        pool,
		broadcaster: broadcaster,
	}
}

func (h *broadcastTransactionHandler) Handle(cmd BroadcastTransaction) error {
	tx, err := h.pool.Get(transaction.ID(cmd.TransactionID))
	if err != nil {
		return err
	}
	if tx == nil {
		return fmt.Errorf("broadcast: transaction not found")
	}

	return h.broadcaster.Broadcast(*tx)
}
