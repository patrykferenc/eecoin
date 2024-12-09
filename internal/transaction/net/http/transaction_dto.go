package http

import "github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"

type inputDTO struct {
	OutputID    string `json:"output_id"`
	OutputIndex int    `json:"output_index"`
	Signature   string `json:"signature"`
}

func (i inputDTO) asInput() *transaction.Input {
	o := transaction.NewInput(transaction.ID(i.OutputID), i.OutputIndex, i.Signature)
	return &o
}

type outputDTO struct {
	Amount  int    `json:"amount"`
	Address string `json:"address"`
}

func (o outputDTO) asOutput() *transaction.Output {
	return transaction.NewOutput(o.Amount, o.Address)
}

type transactionDTO struct {
	ID      string      `json:"id"`
	Inputs  []inputDTO  `json:"inputs"`
	Outputs []outputDTO `json:"outputs"`
}

func asDTO(tx transaction.Transaction) transactionDTO {
	inputs := make([]inputDTO, len(tx.Inputs()))
	for i, in := range tx.Inputs() {
		inputs[i] = inputDTO{
			OutputID:    in.OutputID().String(),
			OutputIndex: in.OutputIndex(),
			Signature:   in.Signature(),
		}
	}

	outputs := make([]outputDTO, len(tx.Outputs()))
	for i, out := range tx.Outputs() {
		outputs[i] = outputDTO{
			Amount:  out.Amount(),
			Address: out.Address(),
		}
	}

	return transactionDTO{
		ID:      tx.ID().String(),
		Inputs:  inputs,
		Outputs: outputs,
	}
}
