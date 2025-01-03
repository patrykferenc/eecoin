package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gymshark/go-hasher"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

var (
	InvalidContentHash    = "Invalid Content Hash"
	BlockNotFound         = errors.New("block not found")
	BlockNotValid         = errors.New("block is not valid")
	BlockDidNotMatchDiff  = errors.New("block did not match difficulty")
	BlockWasNotWithinTime = errors.New("block was not within time")
	ChainNotValid         = errors.New("chain not valid")
	GenesisBlockTimestamp = time.Date(2024, 11, 16, 20, 23, 0, 0, time.UTC).UnixMilli()
)

type Block struct {
	Index          int
	TimestampMilis int64
	ContentHash    string
	PrevHash       string
	Transactions   []transaction.Transaction
	Challenge      Challenge
}

func (block Block) MarshalBinary() ([]byte, error) { // TODO#30 - maybe prettify this
	var buffer bytes.Buffer
	_, err := buffer.WriteString(fmt.Sprintf("%d%d%s", block.Index, block.TimestampMilis, block.PrevHash))
	if err != nil {
		return nil, fmt.Errorf("error writing to buffer: %w", err)
	}

	marshalledChallange, err := block.Challenge.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("error marshalling block: %w", err)
	}
	_, err = buffer.Write(marshalledChallange)
	if err != nil {
		return nil, fmt.Errorf("error writing to buffer: %w", err)
	}

	for _, transaction := range block.Transactions {
		marshalledTransaction, err := transaction.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("error marshalling block: %w", err)
		}
		_, err = buffer.Write(marshalledTransaction)
		if err != nil {
			return nil, fmt.Errorf("error writing to buffer: %w", err)
		}
	}

	return buffer.Bytes(), nil
}

type BlockChain struct {
	Blocks []Block
}

func (chain *BlockChain) NewBlock(timestamp int64, transactions []transaction.Transaction, solved Challenge) (Block, error) {
	if !solved.MatchesDifficulty() {
		slog.Error("Block not valid", "reason", "difficulty not met")
		return Block{}, BlockDidNotMatchDiff
	}
	if !blockCreatedAfterPreviousWithinTimeCap(timestamp, solved, chain.GetLast()) {
		slog.Error("Block not valid", "reason", "time cap not met")
		return Block{}, BlockWasNotWithinTime
	}
	previousHash := chain.Blocks[len(chain.Blocks)-1].ContentHash
	newBlock := &Block{
		Index:          len(chain.Blocks),
		TimestampMilis: timestamp,
		PrevHash:       previousHash,
		Transactions:   transactions,
		Challenge:      solved,
	}
	contentHash, err := CalculateHash(*newBlock)
	if err != nil {
		return Block{}, err
	}
	newBlock.ContentHash = contentHash
	return *newBlock, nil
}

func (chain *BlockChain) AddBlock(new Block) error {
	if isValidBasedOnPrevious(new, chain.GetLast()) {
		chain.Blocks = append(chain.Blocks, new)
		return nil
	}
	return BlockNotValid
}

func (chain *BlockChain) RemoveBlocksStartingWithIndex(index int) {
	shortenedChain := chain.Blocks[:len(chain.Blocks)-index]
	chain.Blocks = shortenedChain
}

func (chain *BlockChain) GetBlock(index int) (Block, error) {
	if index >= len(chain.Blocks) || index < 0 {
		return Block{}, BlockNotFound
	}
	return chain.Blocks[index], nil
}

func (chain *BlockChain) GetLast() Block {
	return chain.Blocks[len(chain.Blocks)-1]
}

func (chain *BlockChain) GetFirst() Block {
	return chain.Blocks[0]
}

func (chain *BlockChain) GetBlockByHash(hash string) (Block, error) {
	for _, block := range chain.Blocks {
		if block.ContentHash == hash {
			return block, nil
		}
	}
	return Block{}, BlockNotFound
}

func (chain *BlockChain) GetBlockByTransactionID(id transaction.ID) (Block, error) {
	for _, block := range chain.Blocks {
		for _, transaction := range block.Transactions {
			if transaction.ID() == id {
				return block, nil
			}
		}
	}
	return Block{}, BlockNotFound
}

func (chain *BlockChain) GetCumulativeDifficulty() int64 {
	var sum int64 = 0
	for _, block := range chain.Blocks {
		sum += int64(intPow(block.Challenge.Difficulty, 2))
	}
	return sum
}

func ImportBlockchain(blocks []Block) (*BlockChain, error) {
	if len(blocks) == 0 || !isValidGenesis(blocks[0]) {
		return nil, ChainNotValid
	}
	for i := 1; i < len(blocks); i++ {
		if !isValidBasedOnPrevious(blocks[i], blocks[i-1]) {
			return nil, ChainNotValid
		}
	}
	return &BlockChain{
		Blocks: blocks,
	}, nil
}

func GenerateGenesisBlock() Block {
	genesisTransaction, _ := transaction.NewGenesis() // todo add error handling
	genesisBlock := &Block{
		Index:          0,
		TimestampMilis: GenesisBlockTimestamp,
		Transactions: []transaction.Transaction{
			*genesisTransaction,
		},
		Challenge: Challenge{
			TimeCapMillis: 1,
			Difficulty:    9,
		},
	}
	contentHash, _ := CalculateHash(*genesisBlock)
	genesisBlock.ContentHash = contentHash
	return *genesisBlock
}

func CalculateHash(block Block) (string, error) {
	blockAsBytes, err := block.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("error calculating hash: %w", err)
	}
	structureHash := hasher.Sha256(blockAsBytes).Base64()
	return structureHash, nil
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
		return newBlock.Index == previous.Index+1 && newBlock.PrevHash == previous.ContentHash &&
			Verify(previous, newBlock.TimestampMilis, newBlock.Challenge.Nonce, newBlock.Challenge.HashValue, newBlock.Transactions) &&
			blockCreatedAfterPreviousWithinTimeCap(newBlock.TimestampMilis, newBlock.Challenge, previous)
	}
	return false
}

func blockCreatedAfterPreviousWithinTimeCap(timestamp int64, solved Challenge, latest Block) bool {
	return timestamp-latest.TimestampMilis >= solved.TimeCapMillis
}

func intPow(n, m int) int {
	if m == 0 {
		return 1
	}

	if m == 1 {
		return n
	}

	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}
