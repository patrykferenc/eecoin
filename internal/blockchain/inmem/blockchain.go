package inmem

import (
	"fmt"
	"github.com/patrykferenc/eecoin/internal/blockchain/inmem/persistence"
	"sync"
	"time"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/common/event"
)

type BlockChain struct {
	chain     *blockchain.BlockChain
	rw        sync.RWMutex
	publisher event.Publisher // TODO#30 - we will refactor this class and send the event from the command handler
}

func NewBlockChain(publisher event.Publisher) (*BlockChain, error) {
	genesis := blockchain.GenerateGenesisBlock()
	ch, err := blockchain.ImportBlockchain([]blockchain.Block{genesis})
	if err != nil {
		return nil, err
	}
	return &BlockChain{
		chain:     ch,
		publisher: publisher,
	}, nil
}

func LoadPersistedBlockchain(path string) (*BlockChain, error) {
	ch, err := persistence.Load(path)
	if err != nil {
		return nil, err
	}
	return &BlockChain{chain: ch}, nil
}

func (b *BlockChain) GetChain() blockchain.BlockChain {
	return *b.chain
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

// TODO add solved challange in order to create a new block
func (b *BlockChain) MarkSeen(id blockchain.TransactionID) error {
	b.rw.Lock()
	// TODO solved challange should be injected in NewBlock here
	block, err := b.chain.NewBlock(time.Now().UnixMilli(), []blockchain.TransactionID{id}, blockchain.Challenge{})
	if err != nil {
		b.rw.Unlock()
		return fmt.Errorf("could not mark block as seen: new block: %w", err)
	}
	b.rw.Unlock()

	e, err := event.New(blockchain.NewBlockAddedEvent{Block: block}, "x.block.added")
	if err != nil {
		return fmt.Errorf("could not mark block as seen: new event: %w", err)
	}

	if err := b.publisher.Publish(e); err != nil {
		return fmt.Errorf("could not mark block as seen: publish event: %w", err)
	}

	return nil
}
