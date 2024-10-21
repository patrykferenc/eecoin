package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node"
	"github.com/patrykferenc/eecoin/internal/peer"
)

type container struct {
	nodeComponent *node.Component
	peerComponent *peer.Component

	broker *event.ChannelBroker
}

func newContainer() (*container, error) {
	file, err := getPeersFile()
	if err != nil {
		return nil, err
	}

	broker := event.NewChannelBroker()

	peerComponent, err := peer.NewComponent(file)
	if err != nil {
		return nil, err
	}

	nodeComponent, err := node.NewComponent(broker, peerComponent.Queries.GetPeers)
	if err != nil {
		return nil, err
	}

	return &container{
		peerComponent: &peerComponent,
		nodeComponent: &nodeComponent,

		broker: broker,
	}, nil
}

func getPeersFile() (io.ReadCloser, error) {
	file, err := os.Open("/etc/eecoin/peers") // TODO: make it configurable
	if err != nil {
		if err != os.ErrNotExist {
			return nil, err
		}
		slog.Warn("Peers file not found, creating a new one")
		file, err = os.Create("/etc/eecoin/peers")
		if err != nil {
			slog.Error("Failed to create peers file", "error", err)
			return nil, err
		}
	}

	return file, nil
}
