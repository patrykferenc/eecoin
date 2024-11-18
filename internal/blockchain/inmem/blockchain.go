package inmem

import (
	"fmt"
	"sync"
	"time"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
)

type BlockChain struct {
	chain *blockchain.BlockChain
	rw    sync.RWMutex
}

func NewBlockChain() (*BlockChain, error) {
	genesis := blockchain.GenerateGenesisBlock()
	ch, err := blockchain.ImportBlockchain([]blockchain.Block{genesis})
	if err != nil {
		return nil, err
	}
	return &BlockChain{
		chain: ch,
	}, nil
}

func (b *BlockChain) Seen(id blockchain.TransactionID) (bool, error) {
	b.rw.RLock()
	defer b.rw.RUnlock()

	_, err := b.chain.GetBlockByTransactionID(id)
	if err != nil {
		if err == blockchain.BlockNotFound {
			return false, nil
		}
		return false, fmt.Errorf("could not check if block is seen: get block by transaction id: %w", err)
	}

	return true, nil
}

func (b *BlockChain) MarkSeen(id blockchain.TransactionID) error {
	block, err := b.chain.NewBlock(time.Now().UnixMilli(), []blockchain.TransactionID{id})
	if err != nil {
		return fmt.Errorf("could not mark block as seen: new block: %w", err)
	}

	b.rw.Lock()
	defer b.rw.Unlock()
	return b.chain.AddBlock(block)
}
