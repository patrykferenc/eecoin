package query

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldGetChain(t *testing.T) {
	genBlock := blockchain.GenerateGenesisBlock()
	chain, _ := blockchain.ImportBlockchain([]blockchain.Block{genBlock})
	repo := mockedSeenTransactionRepository{chain: *chain}

	getChain := NewGetChain(&repo)
	result := getChain.Get()
	assert.Equal(t, *chain, result)
}

type mockedSeenTransactionRepository struct {
	chain blockchain.BlockChain
}

func (m *mockedSeenTransactionRepository) Seen(id blockchain.TransactionID) (bool, error) {
	return true, nil
}
func (m *mockedSeenTransactionRepository) MarkSeen(id blockchain.TransactionID) error {
	return nil

}
func (m *mockedSeenTransactionRepository) GetChain() blockchain.BlockChain {
	return m.chain
}
func (m *mockedSeenTransactionRepository) PutBlock(block blockchain.Block) error {
	return m.chain.AddBlock(block)
}
