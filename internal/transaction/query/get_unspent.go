package query

import "github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"

type GetUnspentOutputs interface {
	Get() (UnspentOutputs, error)
}

type UnspentOutputs struct {
	Outputs []UnspentOutput `json:"outputs"`
	Count   int             `json:"count"`
}

type UnspentOutput struct {
	OutputID    string `json:"output_id"`
	OutputIndex int    `json:"output_index"`
	Amount      int    `json:"amount"`
	Address     string `json:"address"`
}

func (u UnspentOutputs) ToModel() []transaction.UnspentOutput {
	oo := make([]transaction.UnspentOutput, u.Count)
	for i, o := range u.Outputs {
		oo[i] = transaction.NewUnspentOutput(transaction.ID(o.OutputID), o.OutputIndex, o.Amount, o.Address)
	}
	return oo
}

func unspentOutputsFromModel(unspents []transaction.UnspentOutput) UnspentOutputs {
	uu := UnspentOutputs{
		Outputs: make([]UnspentOutput, len(unspents)),
		Count:   len(unspents),
	}
	for i, unspent := range unspents {
		uu.Outputs[i] = UnspentOutput{
			OutputID:    string(unspent.OutputID()),
			OutputIndex: unspent.OutputIndex(),
			Amount:      unspent.Amount(),
			Address:     unspent.Address(),
		}
	}
	return uu
}

type getUnspentOutputs struct {
	repo transaction.UnspentOutputRepository
}

func NewGetUnspentOutputs(repo transaction.UnspentOutputRepository) GetUnspentOutputs {
	return &getUnspentOutputs{repo: repo}
}

func (g *getUnspentOutputs) Get() (UnspentOutputs, error) {
	oo, err := g.repo.GetAll()
	if err != nil {
		return UnspentOutputs{}, err
	}
	return unspentOutputsFromModel(oo), nil
}
