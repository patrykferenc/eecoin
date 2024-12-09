package inmem

import (
	"sync"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type UnspentOutputRepository struct {
	outputs map[string][]transaction.UnspentOutput
	rw      sync.RWMutex
}

func NewUnspentOutputRepository() *UnspentOutputRepository {
	return &UnspentOutputRepository{
		outputs: make(map[string][]transaction.UnspentOutput),
	}
}

func (r *UnspentOutputRepository) GetAll() ([]transaction.UnspentOutput, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	var uos []transaction.UnspentOutput
	for _, outputs := range r.outputs {
		for _, output := range outputs {
			uos = append(uos, output)
		}
	}

	return uos, nil
}

func (r *UnspentOutputRepository) GetByAddress(address string) ([]transaction.UnspentOutput, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	return r.outputs[address], nil
}

func (r *UnspentOutputRepository) Set(outputs []transaction.UnspentOutput) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	for _, output := range outputs {
		r.outputs[output.Address()] = append(r.outputs[output.Address()], output)
	}

	return nil
}
