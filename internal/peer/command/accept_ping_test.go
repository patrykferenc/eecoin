package command_test

import (
	"fmt"
	"testing"

	"github.com/patrykferenc/eecoin/internal/peer/command"
	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
)

func TestPingSendAndAccept(t *testing.T) {
	// given two peers, one one-sided
	peersA := peer.NewPeers([]*peer.Peer{
		{Host: "hostB", Status: peer.StatusUnknown},
		{Host: "hostC", Status: peer.StatusHealthy},
		{Host: "hostD", Status: peer.StatusUnhealthy},
	})
	peersCtxA := &simplePeersContext{peers: peersA}
	peersB := peer.NewPeers([]*peer.Peer{
		{Host: "hostC", Status: peer.StatusHealthy},
	})
	peersCtxB := &simplePeersContext{peers: peersB}

	// and given mocked connection
	config := coordinatedPings{
		errs: map[string]error{
			"hostB": nil,
			"hostC": fmt.Errorf("ping failed"),
			"hostD": fmt.Errorf("ping failed"),
		},
		handlers: map[string]command.AcceptPingHandler{
			"hostA": command.NewAcceptPingHandler(peersCtxA),
			"hostB": command.NewAcceptPingHandler(peersCtxB),
			"hostC": &noOpAcceptPingHandler{},
			"hostD": &noOpAcceptPingHandler{},
		},
	}
	coordinatedPingSenderA := &coordinatedPingSender{config: config, sourceHost: "hostA", t: t}

	// when
	senderA := command.NewSendPingHandler(coordinatedPingSenderA, peersCtxA)
	senderA.Handle(command.SendPingCommand{})

	// then
	for _, p := range peersA.All() {
		if p.Host == "hostB" && p.Status != peer.StatusHealthy {
			t.Errorf("expected peer %s to be healthy, got %s", p.Host, p.Status)
		}
		if p.Host == "hostC" && p.Status != peer.StatusUnhealthy {
			t.Errorf("expected peer %s to be unhealthy, got %s", p.Host, p.Status)
		}
		if p.Host == "hostD" && p.Status != peer.StatusUnhealthy {
			t.Errorf("expected peer %s to be unhealthy, got %s", p.Host, p.Status)
		}
	}

	// and then should update peerA on peer B which haven't seen A
	seenPeerA := false
	for _, p := range peersB.All() {
		if p.Host == "hostA" {
			seenPeerA = true
			if p.Status != peer.StatusHealthy {
				t.Errorf("expected peer %s to be healthy, got %s", p.Host, p.Status)
			}
		}
	}
	if !seenPeerA {
		t.Errorf("expected peer hostA to be present and seen by hostB")
	}
}

type coordinatedPings struct {
	errs     map[string]error // if pinging host should error (simulate connection)
	handlers map[string]command.AcceptPingHandler
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
		err := handler.Handle(command.AcceptPing{Host: s.sourceHost})
		if err != nil {
			return err
		}
	}
	return shouldErr
}

type noOpAcceptPingHandler struct{}

func (h *noOpAcceptPingHandler) Handle(cmd command.AcceptPing) error {
	return nil
}
