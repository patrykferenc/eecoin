package peer

import (
	"io"

	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
	"github.com/patrykferenc/eecoin/internal/peer/query"
)

type Component struct {
	Queries Queries
}

type Queries struct {
	query.GetPeers
}

func NewComponent(peersFile io.ReadCloser) (Component, error) {
	defer peersFile.Close()
	peersFromFile, err := peer.PeersFromFile(peersFile)
	if err != nil {
		return Component{}, err
	}
	context := &inMemoryPeerContext{
		peers: peersFromFile,
	}

	return Component{
		Queries: Queries{
			GetPeers: query.NewGetPeers(context),
		},
	}, nil
}

type inMemoryPeerContext struct {
	peers *peer.Peers
}

func (i *inMemoryPeerContext) Peers() *peer.Peers {
	return i.peers
}
