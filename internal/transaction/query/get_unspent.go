package query

import "github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"

type GetUnspentOutputs interface {
	Get() (UnspentOutputs, error)
}

type UnspentOutputs struct {
	Outputs []UnspentOutput `json:"outputs"`
	Count   int             `json:"count"`
}

func (u UnspentOutputs) ToModel() []transaction.UnspentOutput {
	oo := make([]transaction.UnspentOutput, u.Count)
	for i, o := range u.Outputs {
		oo[i] = transaction.NewUnspentOutput(transaction.ID(o.OutputID), o.OutputIndex, o.Amount, o.Address)
	}
	return oo
}

func unspentOutputsFromModel(oo []transaction.UnspentOutput) UnspentOutputs {
	u := UnspentOutputs{
		Outputs: make([]UnspentOutput, len(oo)),
		Count:   len(oo),
	}
	for i, o := range oo {
		u.Outputs[i] = UnspentOutput{
			OutputID:    string(o.OutputID()),
			OutputIndex: o.OutputIndex(),
			Amount:      o.Amount(),
			Address:     o.Address(),
		}
	}
	return u
}

type UnspentOutput struct {
	OutputID    string `json:"output_id"`
	OutputIndex int    `json:"output_index"`
	Amount      int    `json:"amount"`
	Address     string `json:"address"`
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
		panic(err)
	}
	return unspentOutputsFromModel(oo), nil
}
