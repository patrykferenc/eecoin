package blockchain

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/net/http"
	"github.com/patrykferenc/eecoin/internal/blockchain/query"
	peersquery "github.com/patrykferenc/eecoin/internal/peer/query"
)

type Component struct {
	Queries  Queries
	Commands Commands
}

type Queries struct {
	GetChain query.GetChain
}

type Commands struct {
	AddBlock  command.AddBlockHandler
	Broadcast command.BroadcastBlockHandler
}

func NewComponent(repo command.BlockChainRepository, peers peersquery.GetPeers) Component {
	broadcaster := http.NewBroadcaster(nil)

	broadcastHandler := command.NewBroadcastBlockHandler(repo, broadcaster, peers)
	return Component{
		Queries: Queries{
			GetChain: query.NewGetChain(repo),
		},
		Commands: Commands{
			AddBlock:  command.NewAddBlockHandler(repo),
			Broadcast: broadcastHandler,
		},
	}
}
