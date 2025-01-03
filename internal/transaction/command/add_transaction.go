package command

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

// AddTransaction is a command to add a transaction to the pool
type AddTransaction struct {
	ProvidedID string
	Inputs     []*transaction.Input
	Outputs    []*transaction.Output
}

func (c AddTransaction) toTransaction() (*transaction.Transaction, error) {
	if len(c.Inputs) == 0 {
		return nil, fmt.Errorf("no inputs provided")
	}
	if len(c.Outputs) == 0 {
		return nil, fmt.Errorf("no outputs provided")
	}

	tx, err := transaction.NewFrom(c.Inputs, c.Outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	// TODO#30 - seems unnecessary
	//if c.ProvidedID != tx.ID().String() {
	//return nil, fmt.Errorf("provided ID [%s] does not match transaction ID [%s]", c.ProvidedID, tx.ID().String())
	//}

	return tx, nil
}

type AddTransactionHandler interface {
	Handle(AddTransaction) error
}

type addTransactionHandler struct {
	pool      transaction.Pool
	publisher event.Publisher
}

func NewAddTransactionHandler(
	publisher event.Publisher,
	pool *transaction.Pool,
) AddTransactionHandler {
	return &addTransactionHandler{
		pool:      *pool,
		publisher: publisher,
	}
}

func (h *addTransactionHandler) Handle(c AddTransaction) error {
	tx, err := c.toTransaction()
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	err = h.pool.Add(tx)
	if err != nil {
		return fmt.Errorf("error adding transaction to pool: %w", err)
	}

	event, err := event.New(transaction.Added{ID: transaction.ID(c.ProvidedID)}, "x.transaction.added")
	if err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}
	if err := h.publisher.Publish(event); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}
