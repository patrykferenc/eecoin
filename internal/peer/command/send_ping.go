package command

import (
	"log/slog"

	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
)

type SendPingCommand struct{}

type SendPingHandler interface {
	Handle(cmd SendPingCommand)
}

type sendPingHandler struct {
	sender  peer.PingSender
	peerCtx peer.PeerContext
}

func NewSendPingHandler(sender peer.PingSender, peerCtx peer.PeerContext) *sendPingHandler {
	if sender == nil {
		panic("sender is nil")
	}
	if peerCtx == nil {
		panic("peerCtx is nil")
	}
	return &sendPingHandler{sender: sender, peerCtx: peerCtx}
}

func (h *sendPingHandler) Handle(cmd SendPingCommand) {
	peers := h.peerCtx.Peers()
	allPeers := peers.All()
	slog.Info("Pinging", "peers", len(allPeers))

	// TODO: Can be parallelized
	for _, p := range allPeers {
		err := h.sender.Ping(p.Host)
		if err != nil {
			slog.Debug("Ping to failed, marking as unhealthy", "host", p.Host, "err", err)
			peers.UpdatePeerStatus(p.Host, peer.StatusUnhealthy)
		} else {
			slog.Debug("Ping to succeeded, marking as healthy", "host", p.Host)
			peers.UpdatePeerStatus(p.Host, peer.StatusHealthy)
		}
	}
}
