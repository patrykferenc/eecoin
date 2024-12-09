package transaction

import (
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/peer/query"
	"github.com/patrykferenc/eecoin/internal/transaction/command"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/patrykferenc/eecoin/internal/transaction/net/http"
)

type Component struct {
	Queries  Queries
	Commands Commands
}

type Queries struct{}

type Commands struct {
	BroadcastTransactionHandler command.BroadcastTransactionHandler
	AddTransactionHandler       command.AddTransactionHandler
}

func NewComponent(
	publisher event.Publisher,
	poolRepository transaction.PoolRepository,
	getPeers query.GetPeers,
) Component {
	pool := transaction.NewPool(poolRepository)
	add := command.NewAddTransactionHandler(
		publisher,
		pool,
	)
	broadcaster := http.NewBroadcaster(nil, getPeers)
	broadcast := command.NewBroadcastTransactionHandler(
		pool,
		broadcaster,
	)
	return Component{
		Queries: Queries{},
		Commands: Commands{
			AddTransactionHandler:       add,
			BroadcastTransactionHandler: broadcast,
		},
	}
}
