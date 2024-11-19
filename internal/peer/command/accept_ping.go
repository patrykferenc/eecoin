package command

import (
	"fmt"
	"log/slog"

	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
)

type AcceptPing struct {
	Host string
}

type AcceptPingHandler interface {
	Handle(cmd AcceptPing) error
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

func (h *acceptPingHandler) Handle(cmd AcceptPing) error {
	if cmd.Host == "" {
		return fmt.Errorf("host is empty")
	}
	peers := h.peerCtx.Peers()
	slog.Debug("Accepted ping", "host", cmd.Host)
	peers.UpdatePeerStatus(cmd.Host, peer.StatusHealthy)
	return nil
}
