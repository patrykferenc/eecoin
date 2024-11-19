package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPersistBlockChain(t *testing.T) {
	assertThat := assert.New(t)

	// given
	dir := t.TempDir() + "/chain"
	genesis := GenerateGenesisBlock()
	chain, _ := ImportBlockchain([]Block{genesis})

	// and given
	timestamp := time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli()
	transactions := make([]TransactionID, 0)
	newBlock, _ := chain.NewBlock(timestamp, transactions)
	_ = chain.AddBlock(newBlock)

	// when - then
	saveErr := Persist(*chain, dir)
	assertThat.Nil(saveErr)

	//when - then
	loaded, err := Load(dir)

	assertThat.Nil(err)
	assertThat.NotNil(loaded)
	assertThat.Equal(loaded.GetFirst(), chain.GetFirst())
	assertThat.Equal(loaded.GetLast(), chain.GetLast())
}
