package blockchain

import (
	"encoding/json"
	"os"
)

// TODO#29 - Move persistence logic to a separate package
type BlockDto struct {
	Index          int             `json:"index"`
	TimestampMilis int64           `json:"timestampMilis"`
	ContentHash    uint64          `json:"contentHash"`
	PrevHash       uint64          `json:"prevHash"`
	Transactions   []TransactionID `json:"transactions"`
	Challenge      Challenge       `json:"challenge"`
}

type ChainDto struct {
	Blocks []BlockDto `json:"blocks"`
}

func MapToDto(blockchain BlockChain) ChainDto {
	dtoBlocks := make([]BlockDto, len(blockchain.blocks))
	for i, block := range blockchain.blocks {
		dtoBlocks[i] = BlockDto{
			Index:          block.Index,
			TimestampMilis: block.TimestampMilis,
			ContentHash:    block.ContentHash,
			PrevHash:       block.PrevHash,
			Transactions:   block.Transactions,
			Challenge:      block.Challenge,
		}
	}
	return ChainDto{Blocks: dtoBlocks}
}

func MapToActual(blockchain ChainDto) (BlockChain, error) {
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

func Persist(chain BlockChain, path string) error {
	mappedToDto := MapToDto(chain)
	b, err := json.Marshal(mappedToDto)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

func Load(path string) (*BlockChain, error) {
	persistedContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	blockchainDto := ChainDto{}
	err = json.Unmarshal(persistedContent, &blockchainDto)
	if err != nil {
		return nil, err
	}

	blockchain, err := MapToActual(blockchainDto)
	if err != nil {
		return nil, err
	}
	return &blockchain, nil
}
