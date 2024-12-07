package blockchain

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/net/http"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

type Component struct {
	Queries  Queries
	Commands Commands
}

type Queries struct{}

type Commands struct {
	AddBlock  command.AddBlockHandler
	Broadcast command.BroadcastBlockHandler
}

func NewComponent(repo command.BlockChainRepository, peersRepo node.PeersRepository) Component {
	broadcaster := http.NewBroadcaster(nil)

	broadcastHandler := command.NewBroadcastBlockHandler(repo, broadcaster, peersRepo)
	return Component{
		Queries: Queries{},
		Commands: Commands{
			AddBlock:  command.NewAddBlockHandler(repo),
			Broadcast: broadcastHandler,
		},
	}
}
