package peer

import (
	"io"

	"github.com/patrykferenc/eecoin/internal/peer/command"
	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
	"github.com/patrykferenc/eecoin/internal/peer/net/http"
	"github.com/patrykferenc/eecoin/internal/peer/query"
)

type Component struct {
	Queries  Queries
	Commands Commands
}

type Commands struct {
	SendPing   command.SendPingHandler
	AcceptPing command.AcceptPingHandler
	SavePeers  command.SavePeersCommandHandler
}

type Queries struct {
	GetPeers query.GetPeers
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

	sender := http.NewPingClient(command.NewAcceptPingHandler(context))

	return Component{
		Queries: Queries{
			GetPeers: query.NewGetPeers(context),
		},
		Commands: Commands{
			SendPing:   command.NewSendPingHandler(sender, context),
			AcceptPing: command.NewAcceptPingHandler(context),
			SavePeers:  command.NewSavePeersCommandHandler(context),
		},
	}, nil
}

type inMemoryPeerContext struct {
	peers *peer.Peers
}

func (i *inMemoryPeerContext) Peers() *peer.Peers {
	return i.peers
}
