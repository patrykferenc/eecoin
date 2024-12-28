package transaction

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"strings"
)

const COINBASE_AMOUNT = 100

// ID is the transaction ID, represented as a base64 string
type ID string

func (h ID) String() string {
	return string(h)
}

func newID(ins []*Input, outs []*Output) (ID, error) {
	var sb strings.Builder

	for _, in := range ins {
		sb.WriteString(in.OutputId.String())
		sb.WriteString(fmt.Sprint(in.OutputIdx))
	}

	for _, out := range outs {
		sb.WriteString(fmt.Sprint(out.Amoun))
		sb.WriteString(fmt.Sprint(out.Addr))
	}

	h := sha256.New()
	_, err := h.Write([]byte(sb.String()))
	if err != nil {
		return "", fmt.Errorf("error writing to hash: %w", err)
	}

	return ID(h.Sum(nil)), nil
}

type Transaction struct { // TODO#30 - rename
	Id ID
	In []*Input
	Ou []*Output
}

func NewFrom(inputs []*Input, outputs []*Output) (*Transaction, error) {
	id, err := newID(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction ID: %w", err)
	}

	return &Transaction{
		Id: id,
		In: inputs,
		Ou: outputs,
	}, nil
}

func (t Transaction) ID() ID {
	return t.Id
}

// Inputs() returns immutable slice of In
func (t Transaction) Inputs() []Input {
	ii := make([]Input, len(t.In))
	for i, in := range t.In {
		ii[i] = *in
	}
	return ii
}

// Outputs() returns immutable slice of Ou
func (t Transaction) Outputs() []Output {
	oo := make([]Output, len(t.Ou))
	for i, out := range t.Ou {
		oo[i] = *out
	}
	return oo
}

// TODO#30 - Addr can be taken from signer
func New(receiverAddr string, senderAddr string, amount int, pk crypto.Signer, unspentOutputRepository UnspentOutputRepository) (*Transaction, error) {
	unspentOutputs, err := unspentOutputRepository.GetByAddress(senderAddr)
	if err != nil {
		return nil, fmt.Errorf("error getting unspent Ou: %w", err)
	}

	// TODO#38 - filter unspent Ou already present in the pool
	leftover, included, err := calculateUnspentForAmount(unspentOutputs, amount)
	if err != nil {
		return nil, fmt.Errorf("error calculating unspent Ou: %w", err)
	}

	inputs := make([]*Input, len(included))
	for i, unspentOutput := range included {
		inputs[i] = unspentOutput.AsInput()
	}

	outputs := generateOutputsFor(amount, leftover, senderAddr, receiverAddr)
	tx, err := NewFrom(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	for i, in := range tx.In {
		err := in.sign(pk, tx.Id, included[i])
		if err != nil {
			return nil, fmt.Errorf("error signing input: %w", err)
		}
	}

	return tx, nil
}

func NewGenesis(receiverAddr string, amount int) (*Transaction, error) {
	inputs := []*Input{}
	outputs := []*Output{
		NewOutput(amount, receiverAddr),
	}

	id, err := newID(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction ID: %w", err)
	}

	return &Transaction{
		Id: id,
		In: inputs,
		Ou: outputs,
	}, nil
}

func NewCoinbase(receiverAddr string, blockHeight int) (*Transaction, error) {
	in := NewInput("", blockHeight, "")
	inputs := []*Input{
		&in,
	}
	outputs := []*Output{
		NewOutput(COINBASE_AMOUNT, receiverAddr),
	}

	id, err := newID(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction ID: %w", err)
	}

	return &Transaction{
		Id: id,
		In: inputs,
		Ou: outputs,
	}, nil
}
