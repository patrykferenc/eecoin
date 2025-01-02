package main

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path"

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

	broker             *event.ChannelBroker
	interruptionChanel chan bool
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

	poolRepo := transactioninmem.NewPoolRepository()
	tranasactionComponent := transaction.NewComponent(broker, poolRepo, peerComponent.Queries.GetPeers)
	blockChainComponent := blockchain.NewComponent(cfg.Persistence.SelfKey, seenRepo, peerComponent.Queries.GetPeers, broker, poolRepo)

	return &Container{
		peerComponent:        &peerComponent,
		blockChainComponent:  &blockChainComponent,
		transactionComponent: &tranasactionComponent,

		broker:             broker,
		interruptionChanel: make(chan bool),
	}, nil
}

func getPeersFile(peersFilePath string) (io.ReadCloser, error) {
	file, err := os.Open(peersFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error("Failed to open peers file", "error", err)
			return nil, err
		}
		slog.Warn("Peers file not found, creating a new one under", "path", peersFilePath)
		err = ensureBaseDir(peersFilePath)
		if err != nil {
			slog.Error("Failed to create base dir for peers file", "error", err)
			return nil, err
		}

		file, err = os.Create(peersFilePath)
		if err != nil {
			slog.Error("Failed to create peers file", "error", err)
			return nil, err
		}
	}

	return file, nil
}

func ensureBaseDir(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}
