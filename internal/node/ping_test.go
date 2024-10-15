package node

import (
	"fmt"
	"testing"
)

type simplePingSender struct {
	err error
}

func (s *simplePingSender) Ping(host string) error {
	return s.err
}

type simplePeersContext struct {
	peers *Peers
}

func (c *simplePeersContext) Peers() *Peers {
	return c.peers
}

func TestHandlePingUpdatesStatus(t *testing.T) {
	// given
	peers := NewPeers([]*Peer{
		{Host: "host1", Status: StatusHealthy},
		{Host: "host2", Status: StatusHealthy},
		{Host: "host3", Status: StatusUnhealthy},
	})
	// and given
	peersCtx := &simplePeersContext{peers: peers}
	sender := &simplePingSender{}
	handler := &simplePingHandler{sender: sender, peerCtx: peersCtx}

	// and given failing ping
	sender.err = fmt.Errorf("ping failed")

	// when
	handler.Handle(PingCommand{})

	// then
	for _, peer := range peers.All() {
		if peer.Status != StatusUnhealthy {
			t.Errorf("expected peer %s to be unhealthy, got %s", peer.Host, peer.Status)
		}
	}

	// and then given successful ping
	sender.err = nil

	// when
	handler.Handle(PingCommand{})
	// then
	for _, peer := range peers.All() {
		if peer.Status != StatusHealthy {
			t.Errorf("expected peer %s to be healthy, got %s", peer.Host, peer.Status)
		}
	}
}
