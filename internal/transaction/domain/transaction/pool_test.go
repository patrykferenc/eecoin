package transaction_test

import (
	"testing"

	"github.com/patrykferenc/eecoin/internal/common/mock"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction/transactiontest"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)
	// given
	m := mock.NewPoolRepository()
	pool := transaction.NewPool(m)
	// and given transaction
	tx, err := transactiontest.NewGenesisLike("someAddress", 100)
	assert.NoError(err)

	// when adding a transaction
	err = pool.Add(tx)

	// then transaction should be added
	assert.NoError(err)
	// and exists
	assert.True(pool.Exists(tx.ID()))
}

func TestUpdatePool(t *testing.T) {
	t.Skipf("TODO#30")
}
