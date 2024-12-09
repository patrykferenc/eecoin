package transaction

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"strings"
)

// ID is the transaction ID, represented as a base64 string
type ID string

func (h ID) String() string {
	return string(h)
}

func newID(ins []Input, outs []Output) (ID, error) {
	var sb strings.Builder

	for _, in := range ins {
		sb.WriteString(in.outputID.String())
		sb.WriteString(fmt.Sprint(in.outputIndex))
	}

	for _, out := range outs {
		sb.WriteString(fmt.Sprint(out.amount))
		sb.WriteString(fmt.Sprint(out.address))
	}

	h := sha256.New()
	_, err := h.Write([]byte(sb.String()))
	if err != nil {
		return "", fmt.Errorf("error writing to hash: %w", err)
	}

	return ID(h.Sum(nil)), nil
}

type Transaction struct {
	id      ID
	inputs  []Input
	outputs []Output
}

func newFrom(inputs []Input, outputs []Output) (*Transaction, error) {
	id, err := newID(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction ID: %w", err)
	}

	return &Transaction{
		id:      id,
		inputs:  inputs,
		outputs: outputs,
	}, nil
}

func (t Transaction) ID() ID {
	return t.id
}

// TODO#30 - address can be taken from signer
func New(receiverAddr string, senderAddr string, amount int, pk crypto.Signer, unspentOutputRepository UnspentOutputRepository) (*Transaction, error) {
	unspentOutputs, err := unspentOutputRepository.GetByAddress(senderAddr)
	if err != nil {
		return nil, fmt.Errorf("error getting unspent outputs: %w", err)
	}

	// TODO#38 - filter unspent outputs already present in the pool
	leftover, included, err := calculateUnspentForAmount(unspentOutputs, amount)
	if err != nil {
		return nil, fmt.Errorf("error calculating unspent outputs: %w", err)
	}

	inputs := make([]Input, len(included))
	for i, unspentOutput := range included {
		inputs[i] = unspentOutput.AsInput()
	}

	outputs := generateOutputsFor(amount, leftover, senderAddr, receiverAddr)
	tx, err := newFrom(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	// TODO#30 - sign transaction

	return tx, nil
}

func NewGenesis(receiverAddr string, amount int) (*Transaction, error) {
	inputs := []Input{}
	outputs := []Output{
		NewOutput(amount, receiverAddr),
	}

	id, err := newID(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction ID: %w", err)
	}

	return &Transaction{
		id:      id,
		inputs:  inputs,
		outputs: outputs,
	}, nil
}
