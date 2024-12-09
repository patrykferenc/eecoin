package query

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
)

type GetChain interface {
	Get() blockchain.BlockChain
}

type getChain struct {
	repo command.BlockChainRepository
}

func NewGetChain(repo command.BlockChainRepository) GetChain { // TODO: move the repo somewhere else
	return &getChain{repo: repo}
}

func (g *getChain) Get() blockchain.BlockChain {
	return g.repo.GetChain()
}
