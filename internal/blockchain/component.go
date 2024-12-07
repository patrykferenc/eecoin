package blockchain

import "github.com/patrykferenc/eecoin/internal/blockchain/command"

type Component struct {
	Queries  Queries
	Commands Commands
}

type Queries struct{}

type Commands struct {
	AddBlock command.AddBlockHandler
}

func NewComponent(repo command.BlockChainRepository) Component {
	return Component{
		Queries:  Queries{},
		Commands: Commands{AddBlock: command.NewAddBlockHandler(repo)},
	}
}
