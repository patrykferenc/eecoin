package query

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

type GetChain interface {
	Get() blockchain.BlockChain
}

type getChain struct {
	repo node.SeenTransactionRepository
}

func NewGetChain(repo node.SeenTransactionRepository) GetChain {
	return &getChain{repo: repo}
}
func (g *getChain) Get() blockchain.BlockChain {
	return g.repo.GetChain()
}
