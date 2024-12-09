package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/patrykferenc/eecoin/internal/blockchain"
	"github.com/patrykferenc/eecoin/internal/blockchain/inmem"
	"github.com/patrykferenc/eecoin/internal/common/config"
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/peer"
	"github.com/patrykferenc/eecoin/internal/transaction"
	transactioninmem "github.com/patrykferenc/eecoin/internal/transaction/inmem"
)

type Container struct {
	peerComponent        *peer.Component
	blockChainComponent  *blockchain.Component
	transactionComponent *transaction.Component

	broker *event.ChannelBroker
}

func NewContainer(cfg *config.Config) (*Container, error) {
	file, err := getPeersFile(cfg.Peers.FilePath)
	if err != nil {
		return nil, err
	}

	broker := event.NewChannelBroker()

	peerComponent, err := peer.NewComponent(file)
	if err != nil {
		return nil, err
	}

	var seenRepo *inmem.BlockChain
	seenRepo, err = inmem.LoadPersistedBlockchain(cfg.Persistence.ChainFilePath)
	if err != nil {
		slog.Error("couldn't load persistent blockchain, creating new runtime chain", "error", err.Error())
		seenRepo, err = inmem.NewBlockChain(broker)
		if err != nil {
			return nil, err
		}
	}

	blockChainComponent := blockchain.NewComponent(seenRepo, peerComponent.Queries.GetPeers, broker)

	poolRepo := transactioninmem.NewPoolRepository()
	tranasactionComponent := transaction.NewComponent(broker, poolRepo, peerComponent.Queries.GetPeers)

	return &Container{
		peerComponent:        &peerComponent,
		blockChainComponent:  &blockChainComponent,
		transactionComponent: &tranasactionComponent,

		broker: broker,
	}, nil
}

func getPeersFile(peersFilePath string) (io.ReadCloser, error) {
	file, err := os.Open(peersFilePath)
	if err != nil {
		if err != os.ErrNotExist {
			return nil, err
		}
		slog.Warn("Peers file not found, creating a new one under", "path", peersFilePath)
		file, err = os.Create(peersFilePath)
		if err != nil {
			slog.Error("Failed to create peers file", "error", err)
			return nil, err
		}
	}

	return file, nil
}
