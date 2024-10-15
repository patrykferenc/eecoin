package node

import "log/slog"

type Peers struct {
	peers         map[string]*Peer
	peersStatuses map[Status][]*Peer // todo can be map I suppose
}

func (p *Peers) All() []Peer {
	allPeers := make([]Peer, len(p.peers))
	i := 0
	for _, peer := range p.peers {
		allPeers[i] = *peer
		i++
	}
	return allPeers
}

func (p *Peers) Healthy() []Peer {
	healthyPeers := p.peersStatuses[StatusHealthy]
	copiedPeers := make([]Peer, len(healthyPeers))
	for i, peer := range healthyPeers {
		copiedPeers[i] = *peer
	}
	return copiedPeers
}

func (p *Peers) UpdatePeerStatus(host string, status Status) {
	peer, ok := p.peers[host]
	if !ok {
		slog.Warn("peer %s not found but tried to update status to %s", host, status.String())
		return
	}

	currentStatusPeers := p.peersStatuses[peer.Status]
	seen := false
	for i, pp := range currentStatusPeers {
		if pp.Host == host {
			p.peersStatuses[peer.Status] = append(currentStatusPeers[:i], currentStatusPeers[i+1:]...)
			seen = true
			break
		}
	}
	if !seen {
		slog.Warn("peer %s not found in status %s", host, peer.Status.String())
		return
	}

	peer.Status = status

	p.peersStatuses[status] = append(p.peersStatuses[status], peer)
}

func NewPeers(peers []*Peer) *Peers {
	peerMap := make(map[string]*Peer, len(peers))
	peerStatuses := make(map[Status][]*Peer, 3)

	for _, peer := range peers {
		peerMap[peer.Host] = peer
		peerStatuses[peer.Status] = append(peerStatuses[peer.Status], peer)
	}

	return &Peers{
		peers:         peerMap,
		peersStatuses: peerStatuses,
	}
}
