package node

import (
	"strings"
	"testing"
)

func TestReadFromFile(t *testing.T) {
	// given
	givenFile := "localhost:8080\n192.168.26.2:8081"

	// when
	peers, err := PeersFromFile(strings.NewReader(givenFile))
	if err != nil {
		t.Fatalf("ReadFromFile() err = %v; want nil", err)
	}

	// then
	if len(peers.All()) != 2 {
		t.Errorf("Peers should not be empty")
	}
	if len(peers.Healthy()) != 0 {
		t.Errorf("Healthy peers should be empty")
	}
	// and then contains all peers
	for i, peer := range peers.All() {
		if peer.Host != strings.Split(givenFile, "\n")[i] {
			t.Errorf("Host should be %s, got %s", strings.Split(givenFile, "\n")[i], peer.Host)
		}
		if peer.Status != StatusUnknown {
			t.Errorf("Status should be unhealthy, got %v for %s", peer.Status, peer.String())
		}
	}
}

func TestReadFromFileEmpty(t *testing.T) {
	// when
	_, err := PeersFromFile(strings.NewReader(""))
	// then
	if err != nil {
		t.Fatalf("ReadFromFile() err = %v; want nil", err)
	}
}

func TestReadFromFileFail(t *testing.T) {
	t.Skipf("in the future we can implement validating ip/dns name on startup")
}

func TestSaveToFile(t *testing.T) {
	// given
	peers := NewPeers([]*Peer{
		{Host: "localhost:8080", Status: StatusHealthy},
		{Host: "192.168.21.37:8082", Status: StatusUnknown},
	})

	var b strings.Builder

	// when
	err := PeersToFile(peers, &b)
	if err != nil {
		t.Fatalf("SaveToFile() err = %v; want nil", err)
	}

	// then
	got := strings.Split(b.String(), "\n")
	want := strings.Split("localhost:8080\n192.168.21.37:8082\n", "\n")
	length, content, faultyElem := stringSlicesAreEqualWithoutOrder(got, want)
	if !length {
		t.Errorf("Peers should have length 2, got %d", len(got))
	} else if !content {
		t.Errorf("Peers should be %v, got %v, faulty element: %s", want, got, got[faultyElem])
	}
}

func stringSlicesAreEqualWithoutOrder(a, b []string) (length bool, content bool, faultyElem int) {
	if len(a) != len(b) {
		return false, false, -1
	}

	elemsOfA := make(map[string]bool)
	for _, elem := range a {
		elemsOfA[elem] = true
	}
	for i, elem := range b {
		if !elemsOfA[elem] {
			return true, false, i
		}
	}
	return true, true, -1
}
