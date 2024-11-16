package node

import (
	"errors"
	"github.com/mitchellh/hashstructure/v2"
	"math"
	"time"
)

var (
	BlockNotFound         = errors.New("block not found")
	BlockNotValid         = errors.New("block is not valid")
	ChainNotValid         = errors.New("chain not valid")
	GenesisBlockTimestamp = time.Date(2024, 11, 16, 20, 23, 0, 0, time.UTC).UnixMilli()
)

type Block struct {
	Index          int
	TimestampMilis int64
	ContentHash    uint64
	PrevHash       uint64
	Transactions   []Transaction
}

type BlockChain struct {
	blocks []Block
}

func (chain *BlockChain) NewBlock(timestamp int64, transactions []Transaction) (Block, error) {
	previousHash := chain.blocks[len(chain.blocks)-1].ContentHash
	blockWithoutHash := &Block{
		Index:          len(chain.blocks),
		TimestampMilis: timestamp,
		PrevHash:       previousHash,
		Transactions:   transactions,
	}
	contentHash, err := CalculateHash(*blockWithoutHash)
	if err != nil {
		return Block{}, err
	}
	blockWithoutHash.ContentHash = contentHash
	return *blockWithoutHash, nil
}

func (chain *BlockChain) AddBlock(block Block) error {
	if isValidBasedOnPrevious(block, chain.GetLast()) {
		chain.blocks = append(chain.blocks, block)
		return nil
	}
	return BlockNotValid
}
func (chain *BlockChain) RemoveBlocksStartingWithIndex(index int) {
	shortenedChain := chain.blocks[:len(chain.blocks)-index]
	chain.blocks = shortenedChain
}
func (chain *BlockChain) GetBlock(index int) (Block, error) {
	if index >= len(chain.blocks) || index < 0 {
		return Block{}, BlockNotFound
	}
	return chain.blocks[index], nil
}

func (chain *BlockChain) GetLast() Block {
	return chain.blocks[len(chain.blocks)-1]
}

func (chain *BlockChain) GetFirst() Block {
	return chain.blocks[0]
}

func (chain *BlockChain) GetBlockByHash(hash uint64) (Block, error) {
	for _, block := range chain.blocks {
		if block.ContentHash == hash {
			return block, nil
		}
	}
	return Block{}, BlockNotFound
}

func ImportBlockain(blocks []Block) (*BlockChain, error) {
	if len(blocks) == 0 || !isValidGenesis(blocks[0]) {
		return nil, ChainNotValid
	}
	for i := 1; i < len(blocks); i++ {
		if !isValidBasedOnPrevious(blocks[i], blocks[i-1]) {
			return nil, ChainNotValid
		}
	}
	return &BlockChain{
		blocks: blocks,
	}, nil
}

func GenerateGenesisBlock() Block {
	genesisBlock := &Block{
		Index:          0,
		TimestampMilis: GenesisBlockTimestamp,
		Transactions:   []Transaction{},
	}
	contentHash, _ := CalculateHash(*genesisBlock)
	genesisBlock.ContentHash = contentHash
	return *genesisBlock
}
func CalculateHash(block Block) (uint64, error) {
	blockWithoutHash := Block{
		Index:          block.Index,
		TimestampMilis: block.TimestampMilis,
		PrevHash:       block.PrevHash,
		Transactions:   block.Transactions,
	}
	contentHash, err := hashstructure.Hash(&blockWithoutHash, hashstructure.FormatV2, nil)
	if err != nil {
		return uint64(math.NaN()), err
	}
	return contentHash, nil
}

func isValidGenesis(block Block) bool {
	genesisBlock := GenerateGenesisBlock()
	blockActualHash, err := CalculateHash(block)
	if err == nil &&
		block.Index == genesisBlock.Index &&
		block.TimestampMilis == genesisBlock.TimestampMilis &&
		block.ContentHash == genesisBlock.ContentHash &&
		blockActualHash == genesisBlock.ContentHash {
		return true
	}
	return false
}

func isValidBasedOnPrevious(newBlock Block, previous Block) bool {
	contentHash, _ := CalculateHash(newBlock)
	if contentHash == newBlock.ContentHash {
		return newBlock.Index == previous.Index+1 && newBlock.PrevHash == previous.ContentHash
	}
	return false
}
