package query

import "github.com/patrykferenc/eecoin/internal/peer/domain/peer"

type GetPeers interface {
	Get() []string
}

type getPeersQuery struct {
	repo peer.PeerContext
}

func NewGetPeers(repo peer.PeerContext) GetPeers {
	return &getPeersQuery{repo: repo}
}

func (q *getPeersQuery) Get() []string {
	peers := q.repo.Peers().All() // TODO: can be only healthy one day

	peerAddresses := make([]string, 0, len(peers))
	for i, peer := range peers {
		peerAddresses[i] = peer.Host
	}

	return peerAddresses
}
