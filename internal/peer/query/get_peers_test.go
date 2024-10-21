package query_test

import (
	"testing"

	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
	"github.com/patrykferenc/eecoin/internal/peer/query"
	"github.com/stretchr/testify/assert"
)

func TestShouldGetPeers(t *testing.T) {
	peers := peer.NewPeers([]*peer.Peer{{Host: "localhost:8080", Status: peer.StatusHealthy}})
	query := query.NewGetPeers(&mockedPeerContext{peers: peers})

	peerAddresses, err := query.Get()
	assert.NoError(t, err)
	assert.Equal(t, []string{"localhost:8080"}, peerAddresses)
}

type mockedPeerContext struct {
	peers *peer.Peers
}

func (m *mockedPeerContext) Peers() *peer.Peers {
	return m.peers
}
