package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/patrykferenc/eecoin/internal/peer"
)

type container struct {
	peerComponent *peer.Component
}

func newContainer() (*container, error) {
	file, err := getPeersFile()
	if err != nil {
		return nil, err
	}

	peerComponent, err := peer.NewComponent(file)
	if err != nil {
		return nil, err
	}

	return &container{
		peerComponent: &peerComponent,
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
	}

	return file, nil
}
