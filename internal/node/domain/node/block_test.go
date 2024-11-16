package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestImportBlock(t *testing.T) {
	assertThat := assert.New(t)

	//given
	nonGenesisBlock := Block{
		-1,
		time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli(),
		2137,
		2136,
		make([]Transaction, 0),
	}
	genesisBlock := GenerateGenesisBlock()

	//when
	nonGenesisChain, errorWhichShouldBePresent := ImportBlockain([]Block{nonGenesisBlock})
	genesisChain, errorWhichShouldNotBePresent := ImportBlockain([]Block{genesisBlock})

	//then
	assertThat.Nil(nonGenesisChain)
	assertThat.Equal(errorWhichShouldBePresent, ChainNotValid)
	assertThat.NotNil(genesisChain)
	assertThat.Equal(genesisChain.GetLast(), genesisBlock)
	assertThat.Equal(genesisChain.GetFirst(), genesisBlock)
	assertThat.Nil(errorWhichShouldNotBePresent)
}

func TestNewBlockWithAddBlock(t *testing.T) {
	assertThat := assert.New(t)

	//given
	timestamp := time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli()
	transactions := make([]Transaction, 0)
	genesisBlock := GenerateGenesisBlock()
	chain, _ := ImportBlockain([]Block{genesisBlock})

	//when - then

	invalidIndexBlock := Block{
		12,
		time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli(),
		2137,
		genesisBlock.ContentHash,
		make([]Transaction, 0),
	}

	invalidPrevHashBlock := Block{
		1,
		time.Date(2023, 2, 3, 12, 0, 0, 0, time.UTC).UnixMilli(),
		2137,
		2136,
		make([]Transaction, 0),
	}

	newBlock, errorShouldNotBePresent := chain.NewBlock(timestamp, transactions)
	assertThat.NotNil(newBlock)
	assertThat.Nil(errorShouldNotBePresent)

	_ = chain.AddBlock(invalidIndexBlock)
	expectedInvalid, errInvalidIdx := chain.GetBlock(1)
	assertThat.NotEqual(invalidIndexBlock, chain.GetLast())
	assertThat.NotEqual(invalidIndexBlock, expectedInvalid)
	assertThat.NotEqual(invalidIndexBlock, chain.GetFirst())
	assertThat.Equal(errInvalidIdx, BlockNotFound)

	_ = chain.AddBlock(invalidPrevHashBlock)
	expectedInvalid, errInvalidIdx = chain.GetBlock(1)
	assertThat.NotEqual(invalidIndexBlock, chain.GetLast())
	assertThat.NotEqual(invalidIndexBlock, expectedInvalid)
	assertThat.NotEqual(invalidIndexBlock, chain.GetFirst())
	assertThat.Equal(errInvalidIdx, BlockNotFound)

	_ = chain.AddBlock(newBlock)
	expectedActualBlock, errInvalidIdx := chain.GetBlock(1)
	assertThat.Equal(newBlock, chain.GetLast())
	assertThat.Equal(newBlock, expectedActualBlock)
	assertThat.NotEqual(newBlock, chain.GetFirst())
}
