package node

import (
	"fmt"
	"testing"
)

type simplePingSender struct {
	err error
}

func (s *simplePingSender) Ping(targetHost string) error {
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
	handler := &simpleSendPingHandler{sender: sender, peerCtx: peersCtx}

	// and given failing ping
	sender.err = fmt.Errorf("ping failed")

	// when
	handler.Handle(SendPingCommand{})

	// then
	for _, peer := range peers.All() {
		if peer.Status != StatusUnhealthy {
			t.Errorf("expected peer %s to be unhealthy, got %s", peer.Host, peer.Status)
		}
	}

	// and then given successful ping
	sender.err = nil

	// when
	handler.Handle(SendPingCommand{})
	// then
	for _, peer := range peers.All() {
		if peer.Status != StatusHealthy {
			t.Errorf("expected peer %s to be healthy, got %s", peer.Host, peer.Status)
		}
	}
}

func TestPingSendAndAccept(t *testing.T) {
	// given two peers, one one-sided
	peersA := NewPeers([]*Peer{
		{Host: "hostB", Status: StatusUnknown},
		{Host: "hostC", Status: StatusHealthy},
		{Host: "hostD", Status: StatusUnhealthy},
	})
	peersCtxA := &simplePeersContext{peers: peersA}
	peersB := NewPeers([]*Peer{
		{Host: "hostC", Status: StatusHealthy},
	})
	peersCtxB := &simplePeersContext{peers: peersB}

	// and given mocked connection
	config := coordinatedPings{
		errs: map[string]error{
			"hostB": nil,
			"hostC": fmt.Errorf("ping failed"),
			"hostD": fmt.Errorf("ping failed"),
		},
		handlers: map[string]AcceptPingHandler{
			"hostA": NewAcceptPingHandler(peersCtxA),
			"hostB": NewAcceptPingHandler(peersCtxB),
			"hostC": &noOpAcceptPingHandler{},
			"hostD": &noOpAcceptPingHandler{},
		},
	}
	coordinatedPingSenderA := &coordinatedPingSender{config: config, sourceHost: "hostA", t: t}

	// when
	senderA := NewSendPingHandler(coordinatedPingSenderA, peersCtxA)
	senderA.Handle(SendPingCommand{})

	// then
	for _, peer := range peersA.All() {
		if peer.Host == "hostB" && peer.Status != StatusHealthy {
			t.Errorf("expected peer %s to be healthy, got %s", peer.Host, peer.Status)
		}
		if peer.Host == "hostC" && peer.Status != StatusUnhealthy {
			t.Errorf("expected peer %s to be unhealthy, got %s", peer.Host, peer.Status)
		}
		if peer.Host == "hostD" && peer.Status != StatusUnhealthy {
			t.Errorf("expected peer %s to be unhealthy, got %s", peer.Host, peer.Status)
		}
	}

	// and then should update peerA on peer B which haven't seen A
	seenPeerA := false
	for _, peer := range peersB.All() {
		if peer.Host == "hostA" {
			seenPeerA = true
			if peer.Status != StatusHealthy {
				t.Errorf("expected peer %s to be healthy, got %s", peer.Host, peer.Status)
			}
		}
	}
	if !seenPeerA {
		t.Errorf("expected peer hostA to be present and seen by hostB")
	}
}

type coordinatedPings struct {
	errs     map[string]error // if pinging host should error (simulate connection)
	handlers map[string]AcceptPingHandler
}

type coordinatedPingSender struct {
	config     coordinatedPings
	sourceHost string
	t          *testing.T
}

// Ping sends a ping to targetHost and updates the status of the peer
// based on the configuration.
func (s *coordinatedPingSender) Ping(targetHost string) error {
	s.t.Logf("pinging %s from %s", targetHost, s.sourceHost)
	shouldErr, ok := s.config.errs[targetHost]
	if !ok {
		panic("unexpected host: stopping test")
	}
	handler, ok := s.config.handlers[targetHost]
	if !ok {
		panic("unexpected host: stopping test")
	}
	if shouldErr == nil {
		handler.Handle(AcceptPingCommand{Host: s.sourceHost})
	}
	return shouldErr
}

type noOpAcceptPingHandler struct{}

func (h *noOpAcceptPingHandler) Handle(cmd AcceptPingCommand) {}
