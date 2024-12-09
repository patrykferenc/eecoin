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
	senderAddr, err := x509.MarshalPKIXPublicKey(privateSender.Public())
	assert.NoError(err)

	// and given receiver
	privateReceiver, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // only used to generate a public key
	assert.NoError(err)
	receiverAddr, err := x509.MarshalPKIXPublicKey(privateReceiver.Public())
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
	assert.NotNil(tx)
	assert.NotNil(tx.ID())
	assert.Len(tx.Inputs(), 1)
	// and are signed
	assert.NotEmpty(tx.Inputs()[0].Signature())
}

func TestCreateCoinbase(t *testing.T) {
	assert := assert.New(t)
	// given some receiver
	privateReceiver, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // only used to generate a public key
	assert.NoError(err)
	receiverAddr, err := x509.MarshalPKIXPublicKey(privateReceiver.Public()) // this is the node's miner address

	// when
	tx, err := transaction.NewCoinbase(string(receiverAddr), 1)
	assert.NoError(err)

	// then
	assert.NotNil(tx)
	assert.NotNil(tx.ID())
	assert.Len(tx.Inputs(), 1)
	assert.Len(tx.Outputs(), 1)
	assert.Equal(100, tx.Outputs()[0].Amount())
	assert.Equal(string(receiverAddr), tx.Outputs()[0].Address())
}
