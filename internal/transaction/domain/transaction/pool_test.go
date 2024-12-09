package transaction_test

import (
	"testing"

	"github.com/patrykferenc/eecoin/internal/common/mock"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)
	// given
	m := mock.NewPoolRepository()
	pool := transaction.NewPool(m)
	// and given transaction
	tx, err := transaction.NewGenesis("someAddress", 100)
	assert.NoError(err)

	// when adding a transaction
	err = pool.Add(tx)

	// then transaction should be added
	assert.NoError(err)
	// and exists
	assert.True(pool.Exists(tx.ID()))
}
