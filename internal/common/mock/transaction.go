package mock

import "github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"

type UnspentOutputRepository struct {
	Called         int
	UnspentOutputs map[string][]*transaction.UnspentOutput
}

func (r *UnspentOutputRepository) GetByAddress(address string) ([]*transaction.UnspentOutput, error) {
	r.Called++
	return r.UnspentOutputs[address], nil
}
