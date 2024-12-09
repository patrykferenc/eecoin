package transaction

type Output struct {
	amount  int
	address string // TODO#30 address struct?
}

func NewOutput(amount int, address string) *Output {
	return &Output{
		amount:  amount,
		address: address,
	}
}

func (o Output) Amount() int {
	return o.amount
}

func (o Output) Address() string {
	return o.address
}

func generateOutputsFor(amount int, leftover int, senderAddr string, receiverAddr string) []*Output {
	outputs := []*Output{
		NewOutput(amount, receiverAddr),
	}

	if leftover > 0 {
		outputs = append(outputs, NewOutput(leftover, senderAddr))
	}

	return outputs
}
