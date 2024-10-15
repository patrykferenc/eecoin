package node

type PingSender interface {
	Ping(host string) error
}

type PeerContext interface {
	Peers() *Peers
}

type PingCommand struct{}

type PingHandler interface {
	Handle(cmd PingCommand)
}

type simplePingHandler struct {
	sender  PingSender
	peerCtx PeerContext
}

func NewPingHandler(sender PingSender, peerCtx PeerContext) PingHandler {
	if sender == nil {
		panic("sender is nil")
	}
	if peerCtx == nil {
		panic("peerCtx is nil")
	}
	return &simplePingHandler{sender: sender, peerCtx: peerCtx}
}

func (h *simplePingHandler) Handle(cmd PingCommand) {
	peers := h.peerCtx.Peers()
	for _, peer := range peers.All() {
		err := h.sender.Ping(peer.Host)
		if err != nil {
			peers.UpdatePeerStatus(peer.Host, StatusUnhealthy)
		} else {
			peers.UpdatePeerStatus(peer.Host, StatusHealthy)
		}
	}
}
