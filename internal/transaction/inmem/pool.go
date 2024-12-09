package inmem

import (
	"sync"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type PoolRepository struct {
	transactions map[transaction.ID]*transaction.Transaction
	rw           sync.RWMutex
}

func NewPoolRepository() *PoolRepository {
	return &PoolRepository{
		transactions: make(map[transaction.ID]*transaction.Transaction),
	}
}

func (p *PoolRepository) Add(tx *transaction.Transaction) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	p.transactions[tx.ID()] = tx
	return nil
}

func (p *PoolRepository) Exists(tx transaction.ID) bool {
	p.rw.RLock()
	defer p.rw.RUnlock()

	_, ok := p.transactions[tx]
	return ok
}

func (p *PoolRepository) Get(tx transaction.ID) (*transaction.Transaction, error) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	t, ok := p.transactions[tx]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (p *PoolRepository) Remove(txs ...transaction.ID) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	for _, tx := range txs {
		delete(p.transactions, tx)
	}
	return nil
}

func (p *PoolRepository) GetAll() []transaction.Transaction {
	p.rw.RLock()
	defer p.rw.RUnlock()

	var txs []transaction.Transaction
	for _, tx := range p.transactions {
		txs = append(txs, *tx)
	}
	return txs
}
