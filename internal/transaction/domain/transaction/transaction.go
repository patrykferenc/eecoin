package transaction

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"strings"
)

const (
	COINBASE_AMOUNT = 100
	GENESIS_AMOUNT  = 10000
	GENESIS_ADDRESS = "3059301306072a8648ce3d020106082a8648ce3d03010703420004376119d02e6b95174f1c6af6bdc26c4280036104909fc8025dd3ebf8ed524e5abe265b67c1102edd0204ebdc3ab8556fe979be13a51526cea0d414b133061ec3" // public marshaled x509 PKIX
)

// ID is the transaction ID, represented as a base64 string
type ID string

func (h ID) String() string {
	return string(h)
}

func (h ID) MarshalBinary() ([]byte, error) {
	return []byte(h), nil
}

func newID(ins []*Input, outs []*Output) (ID, error) {
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
	inputs  []*Input
	outputs []*Output
}

// ID() returns the transaction ID
func (t Transaction) ID() ID {
	return t.id
}

// Inputs() returns immutable slice of transactions inputs
func (t Transaction) Inputs() []Input {
	ii := make([]Input, len(t.inputs))
	for i, in := range t.inputs {
		ii[i] = *in
	}
	return ii
}

// Outputs() returns immutable slice of transactions outputs
func (t Transaction) Outputs() []Output {
	oo := make([]Output, len(t.outputs))
	for i, out := range t.outputs {
		oo[i] = *out
	}
	return oo
}

func (t Transaction) MarshalBinary() ([]byte, error) {
	var transactionBytes []byte

	for _, in := range t.inputs {
		inBytes, err := in.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("error marshalling input: %w", err)
		}
		transactionBytes = append(transactionBytes, inBytes...)
	}

	for _, out := range t.outputs {
		outBytes, err := out.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("error marshalling output: %w", err)
		}
		transactionBytes = append(transactionBytes, outBytes...)
	}

	transactionIDBytes, err := t.id.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("error marshalling transaction ID: %w", err)
	}
	transactionBytes = append(transactionBytes, transactionIDBytes...)

	return transactionBytes, nil
}

func NewFrom(inputs []*Input, outputs []*Output) (*Transaction, error) {
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

	for i, in := range tx.inputs {
		err := in.sign(pk, tx.id, included[i])
		if err != nil {
			return nil, fmt.Errorf("error signing input: %w", err)
		}
	}

	return tx, nil
}

func NewGenesis() (*Transaction, error) {
	inputs := []*Input{}
	outputs := []*Output{
		NewOutput(GENESIS_AMOUNT, GENESIS_ADDRESS),
	}

	return NewFrom(inputs, outputs)
}

func NewCoinbase(receiverAddr string, blockHeight int) (*Transaction, error) {
	in := NewInput("", blockHeight, "")
	inputs := []*Input{
		in,
	}
	outputs := []*Output{
		NewOutput(COINBASE_AMOUNT, receiverAddr),
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
