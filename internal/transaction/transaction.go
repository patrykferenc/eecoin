package transaction

import (
	"github.com/patrykferenc/eecoin/internal/common/event"
	peerquery "github.com/patrykferenc/eecoin/internal/peer/query"
	"github.com/patrykferenc/eecoin/internal/transaction/application"
	"github.com/patrykferenc/eecoin/internal/transaction/command"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/patrykferenc/eecoin/internal/transaction/inmem"
	"github.com/patrykferenc/eecoin/internal/transaction/net/http"
	"github.com/patrykferenc/eecoin/internal/transaction/query"
)

type Component struct {
	Queries     Queries
	Commands    Commands
	Application Application
}

type Queries struct {
	GetUnspentOutputs  query.GetUnspentOutputs
	GetTransactionPool query.GetTransactionPool
}

type Commands struct {
	BroadcastTransactionHandler command.BroadcastTransactionHandler
	AddTransactionHandler       command.AddTransactionHandler
}

type Application struct {
	TransactionUpdater *application.TransactionUpdater
}

func NewComponent(
	publisher event.Publisher,
	poolRepository *inmem.PoolRepository,
	getPeers peerquery.GetPeers,
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
	unspent := inmem.NewUnspentOutputRepository()
	getUnspent := query.NewGetUnspentOutputs(unspent)

	unspentClient := http.NewUnspentOutputsRepository("noop")

	poolClient := &http.TransactionPoolClient{}

	updater := application.NewTransactionUpdater(
		poolRepository,
		poolClient,
		unspent,
		unspentClient,
		getPeers,
	)

	return Component{
		Queries: Queries{
			GetUnspentOutputs:  getUnspent,
			GetTransactionPool: poolRepository,
		},
		Commands: Commands{
			AddTransactionHandler:       add,
			BroadcastTransactionHandler: broadcast,
		},
		Application: Application{
			TransactionUpdater: updater,
		},
	}
}
