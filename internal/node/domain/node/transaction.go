package node

import (
	"time"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/wallet"
)

type Transaction struct {
	ID        blockchain.TransactionID
	Content   string
	Timestamp time.Time
	From      wallet.ID
	To        wallet.ID
}

// InFlightTransactionRepository is a repository for transactions that are not yet included in a block.
// These transactons may include transactions received from other peers or clients of a node.
type InFlightTransactionRepository interface {
	Save(transaction *Transaction) error
	Get(id blockchain.TransactionID) (*Transaction, error)
	Discard(id blockchain.TransactionID) error
}

type SimpleInFlightTransactionRepository struct {
	transactions map[blockchain.TransactionID]*Transaction
}

func NewSimpleInFlightTransactionRepository() *SimpleInFlightTransactionRepository {
	return &SimpleInFlightTransactionRepository{
		transactions: make(map[blockchain.TransactionID]*Transaction),
	}
}

func (r *SimpleInFlightTransactionRepository) Save(transaction *Transaction) error {
	r.transactions[transaction.ID] = transaction
	return nil
}

func (r *SimpleInFlightTransactionRepository) Get(id blockchain.TransactionID) (*Transaction, error) {
	if _, ok := r.transactions[id]; !ok {
		return nil, nil
	}
	return r.transactions[id], nil
}

func (r *SimpleInFlightTransactionRepository) Discard(id blockchain.TransactionID) error {
	delete(r.transactions, id)
	return nil
}

type SeenTransactionRepository interface {
	Seen(id blockchain.TransactionID) (bool, error)
	MarkSeen(id blockchain.TransactionID) error
}

type SimpleSeenTransactionRepository struct {
	seen map[blockchain.TransactionID]struct{}
}

// MarkSeen marks a transaction as seen. Exposed for testing.
func (r *SimpleSeenTransactionRepository) MarkSeen(id blockchain.TransactionID) error {
	r.seen[id] = struct{}{}
	return nil
}

func NewSimpleSeenTransactionRepository() *SimpleSeenTransactionRepository {
	return &SimpleSeenTransactionRepository{
		seen: make(map[blockchain.TransactionID]struct{}),
	}
}

func (r *SimpleSeenTransactionRepository) Seen(id blockchain.TransactionID) (bool, error) {
	_, ok := r.seen[id]
	return ok, nil
}
