package transaction

import "testing"

func TestGenerateOutputsWhenLeftover(t *testing.T) {
	// given
	amount := 100
	leftover := 50
	senderAddr := "senderAddress"
	receiverAddr := "receiverAddress"

	// when
	outputs := generateOutputsFor(amount, leftover, senderAddr, receiverAddr)

	// then
	if len(outputs) != 2 {
		t.Error("expected 2 outputs")
	}
	if outputs[0].Amount() != amount {
		t.Error("expected amount to be", amount)
	}
	if outputs[0].Address() != receiverAddr {
		t.Error("expected address to be", receiverAddr)
	}
	if outputs[1].Amount() != leftover {
		t.Error("expected leftover to be", leftover)
	}
	if outputs[1].Address() != senderAddr {
		t.Error("expected address to be", senderAddr)
	}
}

func TestGenerateOutputsWhenNoLeftover(t *testing.T) {
	// given
	amount := 100
	leftover := 0
	senderAddr := "senderAddress"
	receiverAddr := "receiverAddress"

	// when
	outputs := generateOutputsFor(amount, leftover, senderAddr, receiverAddr)

	// then
	if len(outputs) != 1 {
		t.Error("expected 1 output")
	}
	if outputs[0].Amount() != amount {
		t.Error("expected amount to be", amount)
	}
	if outputs[0].Address() != receiverAddr {
		t.Error("expected address to be", receiverAddr)
	}
}
