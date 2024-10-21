package command_test

import (
	"fmt"
	"testing"

	"github.com/patrykferenc/eecoin/internal/peer/command"
	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
)

type simplePingSender struct {
	err error
}

func (s *simplePingSender) Ping(targetHost string) error {
	return s.err
}

type simplePeersContext struct {
	peers *peer.Peers
}

func (c *simplePeersContext) Peers() *peer.Peers {
	return c.peers
}

func TestHandlePingUpdatesStatus(t *testing.T) {
	// given
	peers := peer.NewPeers([]*peer.Peer{
		{Host: "host1", Status: peer.StatusHealthy},
		{Host: "host2", Status: peer.StatusHealthy},
		{Host: "host3", Status: peer.StatusUnhealthy},
	})
	// and given
	peersCtx := &simplePeersContext{peers: peers}
	sender := &simplePingSender{}
	handler := command.NewSendPingHandler(sender, peersCtx)

	// and given failing ping
	sender.err = fmt.Errorf("ping failed")

	// when
	handler.Handle(command.SendPingCommand{})

	// then
	for _, p := range peers.All() {
		if p.Status != peer.StatusUnhealthy {
			t.Errorf("expected peer %s to be unhealthy, got %s", p.Host, p.Status)
		}
	}

	// and then given successful ping
	sender.err = nil

	// when
	handler.Handle(command.SendPingCommand{})
	// then
	for _, p := range peers.All() {
		if p.Status != peer.StatusHealthy {
			t.Errorf("expected peer %s to be healthy, got %s", p.Host, p.Status)
		}
	}
}
