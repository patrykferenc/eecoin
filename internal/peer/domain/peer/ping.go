package peer

import "log/slog"

type PingSender interface {
	Ping(targetHost string) error
}

type PeerContext interface {
	Peers() *Peers
}

type SendPingCommand struct{}

type SendPingHandler interface {
	Handle(cmd SendPingCommand)
}

type AcceptPingCommand struct {
	Host string
}

type AcceptPingHandler interface {
	Handle(cmd AcceptPingCommand)
}

type simpleSendPingHandler struct {
	sender  PingSender
	peerCtx PeerContext
}

func NewSendPingHandler(sender PingSender, peerCtx PeerContext) SendPingHandler {
	if sender == nil {
		panic("sender is nil")
	}
	if peerCtx == nil {
		panic("peerCtx is nil")
	}
	return &simpleSendPingHandler{sender: sender, peerCtx: peerCtx}
}

func (h *simpleSendPingHandler) Handle(cmd SendPingCommand) {
	peers := h.peerCtx.Peers()
	allPeers := peers.All()
	slog.Info("Pinging", "peers", len(allPeers))
	for _, peer := range allPeers {
		err := h.sender.Ping(peer.Host)
		if err != nil {
			slog.Debug("Ping to failed, marking as unhealthy", "host", peer.Host, "err", err)
			peers.UpdatePeerStatus(peer.Host, StatusUnhealthy)
			continue
		} else {
			slog.Debug("Ping to succeeded, marking as healthy", "host", peer.Host)
			peers.UpdatePeerStatus(peer.Host, StatusHealthy)
			continue
		}
	}
}

type simpleAcceptPingHandler struct {
	peerCtx PeerContext
}

func NewAcceptPingHandler(peerCtx PeerContext) AcceptPingHandler {
	if peerCtx == nil {
		panic("peerCtx is nil")
	}
	return &simpleAcceptPingHandler{peerCtx: peerCtx}
}

func (h *simpleAcceptPingHandler) Handle(cmd AcceptPingCommand) {
	peers := h.peerCtx.Peers()
	slog.Info("Accepted ping", "host", cmd.Host)
	peers.UpdatePeerStatus(cmd.Host, StatusHealthy)
}
