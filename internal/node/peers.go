package node

import "log/slog"

type Peers struct {
	peersStatuses map[Status]map[string]*Peer
}

func (p *Peers) All() []Peer {
	var allPeers []Peer
	for _, statusPeers := range p.peersStatuses {
		for _, peer := range statusPeers {
			allPeers = append(allPeers, *peer)
		}
	}
	return allPeers
}

func (p *Peers) Healthy() []Peer {
	healthyPeersMap := p.peersStatuses[StatusHealthy]
	copiedPeers := make([]Peer, len(healthyPeersMap))
	i := 0
	for _, peer := range healthyPeersMap {
		copiedPeers[i] = *peer
		i++
	}
	return copiedPeers
}

func (p *Peers) UpdatePeerStatus(host string, status Status) {
	var peer *Peer
	var found bool

	for _, peersMap := range p.peersStatuses {
		if p, ok := peersMap[host]; ok {
			peer = p
			delete(peersMap, host)
			found = true
			break
		}
	}

	if !found {
		slog.Info("new peer %s with status %s", host, status.String())
		peer = &Peer{Host: host, Status: status}
	}

	peer.Status = status

	if p.peersStatuses[status] == nil {
		p.peersStatuses[status] = make(map[string]*Peer)
	}
	p.peersStatuses[status][host] = peer
}

func NewPeers(peers []*Peer) *Peers {
	peerStatuses := make(map[Status]map[string]*Peer, 3)

	for _, peer := range peers {
		if peerStatuses[peer.Status] == nil {
			peerStatuses[peer.Status] = make(map[string]*Peer)
		}
		peerStatuses[peer.Status][peer.Host] = peer
	}

	return &Peers{
		peersStatuses: peerStatuses,
	}
}
