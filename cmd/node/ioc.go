package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/patrykferenc/eecoin/internal/blockchain/inmem"
	"github.com/patrykferenc/eecoin/internal/common/config"
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node"
	nodedomain "github.com/patrykferenc/eecoin/internal/node/domain/node"
	nodehttp "github.com/patrykferenc/eecoin/internal/node/net/http"
	"github.com/patrykferenc/eecoin/internal/peer"
)

type Container struct {
	nodeComponent *node.Component
	peerComponent *peer.Component

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
		slog.Error("couldn't load persistent blockchain, creating new runtime chain", err)
		seenRepo, err = inmem.NewBlockChain()
		if err != nil {
			return nil, err
		}
	}

	inflightRepo := nodedomain.NewSimpleInFlightTransactionRepository()
	sender := nodehttp.NewSender()
	nodeComponent, err := node.NewComponent(broker, peerComponent.Queries.GetPeers, seenRepo, inflightRepo, sender)
	if err != nil {
		return nil, err
	}

	return &Container{
		peerComponent: &peerComponent,
		nodeComponent: &nodeComponent,

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
