package command

import (
	"fmt"
	"log/slog"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
)

type AddBlock struct {
	ToAdd blockchain.Block
}

type AddBlockHandler interface {
	Handle(AddBlock) error
}

func NewAddBlockHandler(repo BlockChainRepository) AddBlockHandler {
	return &addBlockHandler{
		repo: repo,
	}
}

type addBlockHandler struct {
	repo BlockChainRepository
}

type BlockChainRepository interface { // TODO#30 make not public, refactor to not return the blockchain as a whole (unsafe to read)
	GetChain() blockchain.BlockChain
	PutBlock(block blockchain.Block) error
}

func (h *addBlockHandler) Handle(command AddBlock) error {
	block := command.ToAdd
	chain := h.repo.GetChain()

	err := chain.AddBlock(block)
	if err != nil {
		slog.Warn("could not add block to chain in the handler", "error", err)
		return fmt.Errorf("could not add block to chain: %w", err)
	}

	return nil
}
