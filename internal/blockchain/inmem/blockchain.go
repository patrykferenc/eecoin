package inmem

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/inmem/persistence"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/common/event"
)

type BlockChain struct {
	chain     *blockchain.BlockChain
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

func (b *BlockChain) PutBlock(block blockchain.Block) error {
	return b.chain.AddBlock(block)
}
