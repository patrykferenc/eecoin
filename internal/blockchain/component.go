package blockchain

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/net/http"
	"github.com/patrykferenc/eecoin/internal/blockchain/query"
	peersquery "github.com/patrykferenc/eecoin/internal/peer/query"
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
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
	MineBlock command.MineBlockHandler
}

func NewComponent(repo command.BlockChainRepository, peersRepo node.PeersRepository, publisher event.Publisher) Component {
	broadcaster := http.NewBroadcaster(nil)

	broadcastHandler := command.NewBroadcastBlockHandler(repo, broadcaster, peersRepo)
	mineBlockHandler := command.NewMineBlockHandler(repo, publisher)
	return Component{
		Queries: Queries{
			GetChain: query.NewGetChain(repo),
		},
		Commands: Commands{
			AddBlock:  command.NewAddBlockHandler(repo),
			Broadcast: broadcastHandler,
			MineBlock: mineBlockHandler,
		},
	}
}
