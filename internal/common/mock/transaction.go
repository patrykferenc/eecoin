package mock

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type UnspentOutputRepository struct {
	Called         int
	UnspentOutputs map[string][]transaction.UnspentOutput
}

func (r *UnspentOutputRepository) GetByAddress(address string) ([]transaction.UnspentOutput, error) {
	r.Called++
	return r.UnspentOutputs[address], nil
}

func (r *UnspentOutputRepository) GetAll() ([]transaction.UnspentOutput, error) {
	r.Called++
	var uos []transaction.UnspentOutput
	for _, outputs := range r.UnspentOutputs {
		uos = append(uos, outputs...)
	}
	return uos, nil
}

func (r *UnspentOutputRepository) Set(_ []transaction.UnspentOutput) error {
	panic("not implemented")
}

type PoolRepository struct {
	Called       int
	Transactions map[transaction.ID]*transaction.Transaction
}

func (r *PoolRepository) Add(tx *transaction.Transaction) error {
	r.Called++
	r.Transactions[tx.ID()] = tx
	return nil
}

func (r *PoolRepository) Exists(id transaction.ID) bool {
	r.Called++
	_, ok := r.Transactions[id]
	return ok
}

func (r *PoolRepository) Remove(ids ...transaction.ID) error {
	r.Called++
	for _, id := range ids {
		delete(r.Transactions, id)
	}
	return nil
}

func (r *PoolRepository) GetAll() []transaction.Transaction {
	r.Called++
	var txs []transaction.Transaction
	for _, tx := range r.Transactions {
		txs = append(txs, *tx)
	}
	return txs
}

func (r *PoolRepository) Get(id transaction.ID) (*transaction.Transaction, error) {
	r.Called++
	tx, ok := r.Transactions[id]
	if !ok {
		return nil, fmt.Errorf("transaction with ID %s not found", id)
	}
	return tx, nil
}

func NewPoolRepository() *PoolRepository {
	return &PoolRepository{
		Transactions: make(map[transaction.ID]*transaction.Transaction),
	}
}
