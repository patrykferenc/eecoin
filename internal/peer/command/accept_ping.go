package command

import (
	"log/slog"

	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
)

type AcceptPing struct {
	Host string
}

type AcceptPingHandler interface {
	Handle(cmd AcceptPing)
}

type acceptPingHandler struct {
	peerCtx peer.PeerContext
}

func NewAcceptPingHandler(peerCtx peer.PeerContext) AcceptPingHandler {
	if peerCtx == nil {
		panic("peerCtx is nil")
	}
	return &acceptPingHandler{peerCtx: peerCtx}
}

func (h *acceptPingHandler) Handle(cmd AcceptPing) {
	peers := h.peerCtx.Peers()
	slog.Info("Accepted ping", "host", cmd.Host)
	peers.UpdatePeerStatus(cmd.Host, peer.StatusHealthy)
}
