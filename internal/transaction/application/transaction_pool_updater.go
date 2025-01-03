package application

import (
	"fmt"
	"log/slog"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
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
	poolRetriever    TransactionPoolRetriever
	unspent          transaction.UnspentOutputRepository
	unspentRetriever UnspentOutputRetriever
	peers            query.GetPeers
	bc               BlockChainRepository
}

type BlockChainRepository interface { // TODO#30 make not public, refactor to not return the blockchain as a whole (unsafe to read)
	GetChain() blockchain.BlockChain
}

func NewTransactionUpdater(
	pool UpdatableTransactionPoolRepository,
	poolRetriever TransactionPoolRetriever,
	unspent transaction.UnspentOutputRepository,
	unspentRetriever UnspentOutputRetriever,
	bc BlockChainRepository,
	peers query.GetPeers,
) *TransactionUpdater {
	return &TransactionUpdater{
		pool:             pool,
		poolRetriever:    poolRetriever,
		unspent:          unspent,
		unspentRetriever: unspentRetriever,
		peers:            peers,
		bc:               bc,
	}
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

	transactions, err := u.poolRetriever.Get(peers)
	if err != nil {
		return err
	}
	err = u.pool.Set(transactions)
	if err != nil {
		return err
	}
	slog.Info("transaction pool updated from remote", "poolCount", len(transactions), "unspentCount", len(remoteUnspent))
	return nil
}

func (u *TransactionUpdater) UpdateFromBlockchain() error {
	chain := u.bc.GetChain()

	// update unspent from the blockchain
	unspent := make([]transaction.UnspentOutput, 0)
	for _, block := range chain.Blocks {
		for _, tx := range block.Transactions {
			for i, output := range tx.Outputs() {
				unspent = append(unspent, transaction.NewUnspentOutput(tx.ID(), i, output.Amount(), output.Address()))
			}
		}
	}
	err := u.unspent.Set(unspent)
	if err != nil {
		return fmt.Errorf("error updating unspent from blockchain: %w", err)
	}

	slog.Info("unspent updated from blockchain", "unspentCount", len(unspent))

	return nil
}
