package blockchain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDifficultyLowered(t *testing.T) {
	chain := computeChain(10, 3, BlockGenerationDefaultIntervalMillis*BlockGenerationDefaultIntervalMillis/20+1)
	result, err := GetDifficulty(*chain)
	assert.NoError(t, err)
	assert.Less(t, result, 3)
}

func TestGetDifficultyUpped(t *testing.T) {
	chain := computeChain(10, 3, BlockGenerationDefaultIntervalMillis*DifficultyAdjustmentInterval/20-1)
	result, err := GetDifficulty(*chain)
	assert.NoError(t, err)
	assert.Greater(t, result, 3)
}

func computeChain(elements, fixedDifficulty int, fixedTimecapOfChallenges int64) *BlockChain {
	genesisBlock := GenerateGenesisBlock()
	timestamp := genesisBlock.TimestampMilis
	chain, _ := ImportBlockchain([]Block{genesisBlock})
	for i := 1; i < elements+1; i++ {
		challenge, _ := NewChallenge(fixedDifficulty, fixedTimecapOfChallenges)
		_ = challenge.RollUntilMatchesDifficulty(chain.GetLast(), timestamp+int64(i)*fixedTimecapOfChallenges)
		block, _ := chain.NewBlock(timestamp+int64(i)*fixedTimecapOfChallenges, []TransactionID{}, challenge)
		err := chain.AddBlock(block)
		if err != nil {

		}
	}
	return chain
}
