package peer

import "testing"

func TestPeersEmptyWhenCreated(t *testing.T) {
	// when
	actual := NewPeers(nil)

	// then
	if len(actual.All()) != 0 {
		t.Errorf("Peers should be empty when created")
	}
}

func TestPeersNonEmpty(t *testing.T) {
	// given
	peers := []*Peer{
		{Host: "http://localhost:8080", Status: StatusHealthy},
		{Host: "http://localhost:8081", Status: StatusUnhealthy},
		{Host: "http://192.168.21.37:8082", Status: StatusUnknown},
	}

	// when
	actual := NewPeers(peers)

	// then
	if len(actual.All()) != 3 {
		t.Errorf("Peers should not be empty")
	}
	// and then contains all peers
	hosts := map[string]bool{
		"http://localhost:8080":     false,
		"http://localhost:8081":     false,
		"http://192.168.21.37:8082": false,
	}
	for _, peer := range actual.All() {
		if _, ok := hosts[peer.Host]; !ok {
			t.Errorf("Unexpected peer %s", peer.String())
		}
		hosts[peer.Host] = true
	}
	for host, found := range hosts {
		if !found {
			t.Errorf("Expected peer %s not found", host)
		}
	}

	// and then
	if len(actual.Healthy()) != 1 {
		t.Errorf("Healthy peers should not be empty")
	}
	if healthy := actual.Healthy()[0]; healthy.Host != "http://localhost:8080" {
		t.Errorf("Healthy peer should be http://localhost:8080, got %s", healthy.Host)
	}
}

func TestPeersUpdates(t *testing.T) {
	// given
	peers := []*Peer{
		{Host: "http://localhost:8080", Status: StatusHealthy},
		{Host: "http://localhost:8081", Status: StatusUnhealthy},
	}

	// when
	actual := NewPeers(peers)

	// then
	if len(actual.Healthy()) != 1 {
		t.Errorf("Healthy peers should not be empty")
	}

	// and when
	actual.Healthy()[0].Status = StatusUnhealthy

	// then
	if len(actual.Healthy()) != 1 {
		t.Errorf("Healthy peers should not be empty, should not be updated")
	}

	// and when
	actual.UpdatePeerStatus("http://localhost:8080", StatusUnhealthy)

	// then
	if len(actual.Healthy()) != 0 {
		t.Errorf("Healthy peers should be empty, updated, but got %v", actual.Healthy())
	}
}
