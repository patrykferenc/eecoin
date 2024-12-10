package transaction

type Output struct {
	Amoun int
	Addr  string // TODO#30 Addr struct?
}

func NewOutput(amount int, address string) *Output {
	return &Output{
		Amoun: amount,
		Addr:  address,
	}
}

func (o Output) Amount() int {
	return o.Amoun
}

func (o Output) Address() string {
	return o.Addr
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
