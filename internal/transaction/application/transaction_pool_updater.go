package application

import (
	"log/slog"

	"github.com/patrykferenc/eecoin/internal/peer/query"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type TransactionPoolRetriever interface {
	Get(peers []string) ([]transaction.Transaction, error)
}

type UpdatableTransactionPoolRepository interface {
	transaction.PoolRepository
	Set(transactions []transaction.Transaction) error
}

type UnspentOutputRetriever interface {
	Get(peers []string) ([]transaction.UnspentOutput, error)
}

type TransactionUpdater struct {
	pool             UpdatableTransactionPoolRepository
	retriever        TransactionPoolRetriever
	unspent          transaction.UnspentOutputRepository
	unspentRetriever UnspentOutputRetriever
	peers            query.GetPeers
}

func NewTransactionUpdater(
	pool UpdatableTransactionPoolRepository,
	poolRetriever TransactionPoolRetriever,
	unspent transaction.UnspentOutputRepository,
	unspentRetriever UnspentOutputRetriever,
	peers query.GetPeers,
) *TransactionUpdater {
	return &TransactionUpdater{pool: pool, retriever: poolRetriever, unspent: unspent, unspentRetriever: unspentRetriever, peers: peers}
}

func (u *TransactionUpdater) UpdateFromRemote() error {
	peers, err := u.peers.Get()
	if err != nil {
		return err
	}

	remoteUnspent, err := u.unspentRetriever.Get(peers)
	if err != nil {
		return err
	}
	err = u.unspent.Set(remoteUnspent)
	if err != nil {
		return err
	}

	transactions, err := u.retriever.Get(peers)
	if err != nil {
		return err
	}
	err = u.pool.Set(transactions)
	if err != nil {
		return err
	}
	slog.Info("transaction pool updated", "poolCount", len(transactions), "unspentCount", len(remoteUnspent))
	return nil
}
