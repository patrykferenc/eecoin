package blockchain

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/net/http"
	"github.com/patrykferenc/eecoin/internal/blockchain/query"
	"github.com/patrykferenc/eecoin/internal/common/event"
	peersquery "github.com/patrykferenc/eecoin/internal/peer/query"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
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

func NewComponent(selfAddress string, repo command.BlockChainRepository, peers peersquery.GetPeers, publisher event.Publisher, repository transaction.PoolRepository) Component {
	broadcaster := http.NewBroadcaster(nil)

	broadcastHandler := command.NewBroadcastBlockHandler(repo, broadcaster, peers)
	mineBlockHandler := command.NewMineBlockHandler(selfAddress, repo, publisher, repository)
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
