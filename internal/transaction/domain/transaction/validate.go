package transaction

import "fmt"

func validateCoinbase(tx *Transaction, blockHeight int) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil, the first transaction in a block must be a coinbase transaction")
	}
	if ins := len(tx.inputs); ins != 1 {
		return fmt.Errorf("coinbase transaction must have one input, got %d", ins)
	}
	if tx.inputs[0].OutputIndex() != blockHeight {
		return fmt.Errorf("coinbase transaction input must have output index equal to block height, got %d", tx.inputs[0].OutputIndex())
	}
	if outs := len(tx.outputs); outs != 1 {
		return fmt.Errorf("coinbase transaction must have one output, got %d", outs)
	}
	if tx.outputs[0].amount != COINBASE_AMOUNT {
		return fmt.Errorf("coinbase transaction output amount must be %d, got %d", COINBASE_AMOUNT, tx.outputs[0].amount)
	}

	return nil
}
