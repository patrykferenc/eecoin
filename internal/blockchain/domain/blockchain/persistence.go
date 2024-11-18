package blockchain

import (
	"encoding/json"
	"os"
)

type blockDto struct {
	Index          int             `json:"index"`
	TimestampMilis int64           `json:"timestampMilis"`
	ContentHash    uint64          `json:"contentHash"`
	PrevHash       uint64          `json:"prevHash"`
	Transactions   []TransactionID `json:"transactions"`
	Challenge      Challenge       `json:"challenge"`
}

type blockChainDto struct {
	Blocks []blockDto `json:"blocks"`
}

func mapToDto(blockchain BlockChain) blockChainDto {
	dtoBlocks := make([]blockDto, len(blockchain.blocks))
	for i, block := range blockchain.blocks {
		dtoBlocks[i] = blockDto{
			Index:          block.Index,
			TimestampMilis: block.TimestampMilis,
			ContentHash:    block.ContentHash,
			PrevHash:       block.PrevHash,
			Transactions:   block.Transactions,
			Challenge:      block.Challenge,
		}
	}
	return blockChainDto{Blocks: dtoBlocks}
}

func mapToActual(blockchain blockChainDto) (BlockChain, error) {
	dtoBlocks := make([]Block, len(blockchain.Blocks))
	for i, block := range blockchain.Blocks {
		dtoBlocks[i] = Block{
			Index:          block.Index,
			TimestampMilis: block.TimestampMilis,
			ContentHash:    block.ContentHash,
			PrevHash:       block.PrevHash,
			Transactions:   block.Transactions,
			Challenge:      block.Challenge,
		}
	}
	chain, err := ImportBlockchain(dtoBlocks)
	if err != nil {
		return BlockChain{}, err
	}
	return *chain, nil
}

func persist(chain BlockChain, path string) error {
	mappedToDto := mapToDto(chain)
	b, err := json.Marshal(mappedToDto)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

func load(path string) (BlockChain, error) {
	persistedContent, err := os.ReadFile(path)
	if err != nil {
		return BlockChain{}, err
	}

	blockchainDto := blockChainDto{}
	err = json.Unmarshal(persistedContent, &blockchainDto)
	if err != nil {
		return BlockChain{}, err
	}

	blockchain, err := mapToActual(blockchainDto)
	if err != nil {
		return BlockChain{}, err
	}
	return blockchain, nil
}
