package blockchain

import (
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewForkMedMediator(t *testing.T) {
	assertThat := assert.New(t)

	// given
	genesis := GenerateGenesisBlock()
	chain, err := ImportBlockchain([]Block{genesis})
	assertThat.Nil(err)

	// and given
	mediator := NewForkMediator(chain)
	assertThat.NotNil(mediator)
}

func TestForkMediator_AddBlock(t *testing.T) {
	assertThat := assert.New(t)

	// given
	genesis := GenerateGenesisBlock()
	primaryChain, err := ImportBlockchain([]Block{genesis})
	assertThat.Nil(err)

	// and given
	mediator := NewForkMediator(primaryChain)
	assertThat.NotNil(mediator)

	// and given
	timestamp := time.Date(2024, 11, 16, 20, 23, 0, 0, time.UTC).UnixMilli()
	transactions := make([]transaction.Transaction, 0)

	// and given
	solvedChallenge, err := NewChallenge(2, 400000)
	assertThat.Nil(err)
	timestamp = timestamp + 400000

	err = solvedChallenge.RollUntilMatchesDifficulty(genesis, transactions, timestamp)
	assertThat.Nil(err)
	block, err := primaryChain.NewBlock(timestamp, transactions, solvedChallenge)
	assertThat.Nil(err)
	err = mediator.AddBlock(block)
	assertThat.Nil(err)

	// and given another block of the same difficulty which would indicate a fork
	forkChallenge, err := NewChallenge(2, 400000)
	assertThat.Nil(err)
	err = forkChallenge.RollUntilMatchesDifficulty(genesis, transactions, timestamp)
	assertThat.Nil(err)

	secondaryChain := primaryChain.GetShortenedUpToIndex(block.Index)
	forkBlock, err := secondaryChain.NewBlock(timestamp, transactions, forkChallenge)

	assertThat.Nil(err)
	err = mediator.AddBlock(forkBlock)
	assertThat.Nil(err)

	highestValueChain, err := mediator.GetPrimaryChain()
	assertThat.Nil(err)
	assertThat.Equal(primaryChain, highestValueChain)

}
