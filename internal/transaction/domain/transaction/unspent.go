package transaction

import "fmt"

// UnspentOutput represents an unspent output
type UnspentOutput struct {
	outputID    ID
	outputIndex int
	amount      int
	address     string
}

func (o UnspentOutput) OutputID() ID {
	return o.outputID
}

func (o UnspentOutput) Address() string {
	return o.address
}

func (o UnspentOutput) Amount() int {
	return o.amount
}

func (o UnspentOutput) OutputIndex() int {
	return o.outputIndex
}

func NewUnspentOutput(outputID ID, outputIndex int, amount int, address string) UnspentOutput {
	return UnspentOutput{
		outputID:    outputID,
		outputIndex: outputIndex,
		amount:      amount,
		address:     address,
	}
}

func (o UnspentOutput) AsInput() *Input {
	return &Input{
		outputID:    o.outputID,
		outputIndex: o.outputIndex,
	}
}

type UnspentOutputRepository interface {
	GetAll() ([]UnspentOutput, error)
	GetByAddress(address string) ([]UnspentOutput, error)
	GetByOutputIDAndIndex(outputID ID, outputIndex int) (UnspentOutput, error)
	Set(unspentOutputs []UnspentOutput) error
}

func calculateUnspentForAmount(unspentOutputs []UnspentOutput, amount int) (leftover int, included []UnspentOutput, err error) {
	currentAmount := 0
	for _, unspentOutput := range unspentOutputs {
		if currentAmount >= amount {
			break
		}

		currentAmount += unspentOutput.amount
		included = append(included, unspentOutput)
	}

	if currentAmount < amount {
		return 0, nil, fmt.Errorf("not enough unspent Ou to cover the Amoun, have %d, need %d", currentAmount, amount)
	}

	return currentAmount - amount, included, nil
}
