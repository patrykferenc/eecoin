package transaction

import "fmt"

// UnspentOutput represents an unspent output
type UnspentOutput struct {
	outputID    ID
	outputIndex int
	amount      int
	address     string // TODO#30
}

func NewUnspentOutput(outputID ID, outputIndex int, amount int, address string) *UnspentOutput {
	return &UnspentOutput{
		outputID:    outputID,
		outputIndex: outputIndex,
		amount:      amount,
		address:     address,
	}
}

func (o *UnspentOutput) AsInput() *Input {
	return &Input{
		outputID:    o.outputID,
		outputIndex: o.outputIndex,
	}
}

type UnspentOutputRepository interface {
	GetByAddress(address string) ([]*UnspentOutput, error)
}

func calculateUnspentForAmount(unspentOutputs []*UnspentOutput, amount int) (leftover int, included []*UnspentOutput, err error) {
	currentAmount := 0
	for _, unspentOutput := range unspentOutputs {
		if currentAmount >= amount {
			break
		}

		currentAmount += unspentOutput.amount
		included = append(included, unspentOutput)
	}

	if currentAmount < amount {
		return 0, nil, fmt.Errorf("not enough unspent outputs to cover the amount, have %d, need %d", currentAmount, amount)
	}

	return currentAmount - amount, included, nil
}
