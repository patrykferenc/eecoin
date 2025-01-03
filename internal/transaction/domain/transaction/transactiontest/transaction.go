package transactiontest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"

	"github.com/patrykferenc/eecoin/internal/common/mock"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

func NewTransaction() (*transaction.Transaction, error) {
	privateSender, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	senderAddrRaw, err := x509.MarshalPKIXPublicKey(privateSender.Public())
	if err != nil {
		return nil, err
	}
	senderAddr := hex.EncodeToString(senderAddrRaw)

	privateReceiver, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // only used to generate a public key
	if err != nil {
		return nil, err
	}
	receiverAddrRaw, err := x509.MarshalPKIXPublicKey(privateReceiver.Public())
	if err != nil {
		return nil, err
	}
	receiverAddr := hex.EncodeToString(receiverAddrRaw)

	someTransaction, err := NewGenesisLike(string(senderAddr), 100)
	if err != nil {
		return nil, err
	}
	unspentOutputs := map[string][]transaction.UnspentOutput{
		string(senderAddr): {transaction.NewUnspentOutput(someTransaction.ID(), 0, 100, string(senderAddr))},
	}
	unspentOutputRepo := &mock.UnspentOutputRepository{
		UnspentOutputs: unspentOutputs,
	}

	amount := 100

	return transaction.New(string(receiverAddr), string(senderAddr), amount, privateSender, unspentOutputRepo)
}

func NewGenesisLike(receiverAddr string, amount int) (*transaction.Transaction, error) {
	inputs := []*transaction.Input{}
	outputs := []*transaction.Output{
		transaction.NewOutput(amount, receiverAddr),
	}

	return transaction.NewFrom(inputs, outputs)
}
