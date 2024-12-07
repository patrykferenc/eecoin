package node

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
)

type Transaction struct { // TODO#30 - will need to be refactored
	ID        blockchain.TransactionID `json:"id"`
	Content   string                   `json:"content"`
	Timestamp time.Time                `json:"timestamp"`
	From      wallet.ID                `json:"from"`
	To        wallet.ID                `json:"to"`
}

// InFlightTransactionRepository is a repository for transactions that are not yet included in a block.
// These transactons may include transactions received from other peers or clients of a node.
type InFlightTransactionRepository interface {
	Save(transaction *Transaction) error
	Get(id blockchain.TransactionID) (*Transaction, error)
	Discard(id blockchain.TransactionID) error
}

type SimpleInFlightTransactionRepository struct {
	rw           sync.RWMutex
	transactions map[blockchain.TransactionID]*Transaction
}

func NewSimpleInFlightTransactionRepository() *SimpleInFlightTransactionRepository {
	return &SimpleInFlightTransactionRepository{
		transactions: make(map[blockchain.TransactionID]*Transaction),
	}
}

func (r *SimpleInFlightTransactionRepository) Save(transaction *Transaction) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.transactions[transaction.ID] = transaction
	return nil
}

func (r *SimpleInFlightTransactionRepository) Get(id blockchain.TransactionID) (*Transaction, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	r.rw.RLock()
	defer r.rw.RUnlock()
	if _, ok := r.transactions[id]; !ok {
		slog.Warn("transaction not found", "id", id)
		return nil, nil
	}
	return r.transactions[id], nil
}

func (r *SimpleInFlightTransactionRepository) Discard(id blockchain.TransactionID) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	delete(r.transactions, id)
	return nil
}

type SeenTransactionRepository interface {
	Seen(id blockchain.TransactionID) (bool, error)
	MarkSeen(id blockchain.TransactionID) error
	GetChain() blockchain.BlockChain
}

type SimpleSeenTransactionRepository struct {
	rw   sync.RWMutex
	seen map[blockchain.TransactionID]struct{}
}

// MarkSeen marks a transaction as seen. Exposed for testing.
func (r *SimpleSeenTransactionRepository) MarkSeen(id blockchain.TransactionID) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.seen[id] = struct{}{}
	return nil
}

func NewSimpleSeenTransactionRepository() *SimpleSeenTransactionRepository {
	return &SimpleSeenTransactionRepository{
		seen: make(map[blockchain.TransactionID]struct{}),
	}
}

func (r *SimpleSeenTransactionRepository) Seen(id blockchain.TransactionID) (bool, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	_, ok := r.seen[id]
	return ok, nil
}

func (r *SimpleSeenTransactionRepository) GetChain() blockchain.BlockChain {
	return blockchain.BlockChain{}
}
