package persistence

import (
	"encoding/json"
	bc "github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"os"
)

type BlockDto struct {
	Index          int                `json:"index"`
	TimestampMilis int64              `json:"timestampMilis"`
	ContentHash    string             `json:"contentHash"`
	PrevHash       string             `json:"prevHash"`
	Transactions   []bc.TransactionID `json:"transactions"`
	Challenge      bc.Challenge       `json:"challenge"`
}

type ChainDto struct {
	Blocks []BlockDto `json:"blocks"`
}

func MapToDto(blockchain bc.BlockChain) ChainDto {
	dtoBlocks := make([]BlockDto, len(blockchain.Blocks))
	for i, block := range blockchain.Blocks {
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

func MapToActual(chain ChainDto) (bc.BlockChain, error) {
	dtoBlocks := make([]bc.Block, len(chain.Blocks))
	for i, block := range chain.Blocks {
		dtoBlocks[i] = bc.Block{
			Index:          block.Index,
			TimestampMilis: block.TimestampMilis,
			ContentHash:    block.ContentHash,
			PrevHash:       block.PrevHash,
			Transactions:   block.Transactions,
			Challenge:      block.Challenge,
		}
	}
	output, err := bc.ImportBlockchain(dtoBlocks)
	if err != nil {
		return bc.BlockChain{}, err
	}
	return *output, nil
}

func Persist(chain bc.BlockChain, path string) error {
	mappedToDto := MapToDto(chain)
	b, err := json.Marshal(mappedToDto)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

func Load(path string) (*bc.BlockChain, error) {
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
