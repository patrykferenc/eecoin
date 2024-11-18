package blockchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImportBlock_shouldError(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	tt := []struct {
		description string
		chain       []Block
		expectedErr error
	}{
		{
			description: "Non genesis block",
			chain: []Block{
				{
					Index:          -1,
					TimestampMilis: time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli(),
					ContentHash:    2137,
					PrevHash:       2136,
					Transactions:   make([]TransactionID, 0),
					Challenge:      Challenge{},
				},
			},
			expectedErr: ChainNotValid,
		},
	}

	// when
	for _, tc := range tt {
		_, err := ImportBlockchain(tc.chain)

		// then
		assertThat.Equal(tc.expectedErr, err)
	}
}

func TestImportBlock_shouldWork(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	genesisBlock := GenerateGenesisBlock()

	// when
	genesisChain, errorWhichShouldNotBePresent := ImportBlockchain([]Block{genesisBlock})
	// then
	assertThat.Nil(errorWhichShouldNotBePresent)

	// and then
	assertThat.NotNil(genesisChain)
	// and then
	assertThat.Equal(genesisChain.GetLast(), genesisBlock)
	assertThat.Equal(genesisChain.GetFirst(), genesisBlock)
}

func TestNewBlock(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	genesis := GenerateGenesisBlock()
	chain, err := ImportBlockchain([]Block{genesis})
	assertThat.Nil(err)
	// and given
	timestamp := time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli()
	transactions := make([]TransactionID, 0)

	// when
	newBlock, err := chain.NewBlock(timestamp, transactions)

	// then
	assertThat.Nil(err)
	assertThat.NotNil(newBlock)
}

func TestAddBlock_shouldWork(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	genesis := GenerateGenesisBlock()
	chain, err := ImportBlockchain([]Block{genesis})
	assertThat.Nil(err)
	// and given
	timestamp := time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli()
	transactions := make([]TransactionID, 0)

	// and given new block
	newBlock, err := chain.NewBlock(timestamp, transactions)
	assertThat.Nil(err)

	// then
	err = chain.AddBlock(newBlock)
	assertThat.Nil(err)
	// and then
	assertThat.Equal(newBlock, chain.GetLast())
	// and then
	expectedIndex := 1
	actual, err := chain.GetBlock(expectedIndex)
	assertThat.Nil(err)
	assertThat.Equal(newBlock, actual)
}

func TestNewBlockWithAddBlock_shouldNotWork(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	genesis := GenerateGenesisBlock()
	chain, err := ImportBlockchain([]Block{genesis})
	assertThat.Nil(err)

	tt := []struct {
		description string
		block       Block
	}{
		{
			description: "Invalid index",
			block: Block{
				Index:          12,
				TimestampMilis: time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli(),
				ContentHash:    2137,
				PrevHash:       genesis.ContentHash,
				Transactions:   make([]TransactionID, 0),
				Challenge:      Challenge{},
			},
		},
		{
			description: "Invalid prev hash",
			block: Block{
				Index:          1,
				TimestampMilis: time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli(),
				ContentHash:    2137,
				PrevHash:       2136,
				Transactions:   make([]TransactionID, 0),
				Challenge:      Challenge{},
			},
		},
	}

	for _, tc := range tt {
		// when
		err := chain.AddBlock(tc.block)

		// then
		assertThat.NotNil(err)
		// and then
		assertThat.NotEqual(tc.block, chain.GetLast())
		// and then
		_, err = chain.GetBlock(1)
		assertThat.NotNil(err)
		assertThat.Equal(BlockNotFound, err)
	}
}
