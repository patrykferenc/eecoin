package transaction_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"testing"

	"github.com/patrykferenc/eecoin/internal/common/mock"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/stretchr/testify/assert"
)

func TestCreatingTransaction(t *testing.T) {
	assert := assert.New(t)
	// given sender
	privateSender, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(err)
	senderAddr, err := x509.MarshalECPrivateKey(privateSender)
	assert.NoError(err)

	// and given receiver
	privateReceiver, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // only used to generate a public key
	assert.NoError(err)
	receiverAddr, err := x509.MarshalECPrivateKey(privateReceiver)
	assert.NoError(err)

	// and given unspent outputs
	someTransaction, err := transaction.NewGenesis(string(senderAddr), 100)
	assert.NoError(err)
	unspentOutputs := map[string][]transaction.UnspentOutput{
		string(senderAddr): {transaction.NewUnspentOutput(someTransaction.ID(), 0, 100, string(senderAddr))},
	}
	unspentOutputRepo := &mock.UnspentOutputRepository{
		UnspentOutputs: unspentOutputs,
	}

	// and given amount
	amount := 100

	// when creating a transaction
	tx, err := transaction.New(string(receiverAddr), string(senderAddr), amount, privateSender, unspentOutputRepo)
	assert.NoError(err)

	// then transaction should be created
	assert.NotNil(tx) // TOOD#30 assertions
}