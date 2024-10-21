package command

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/patrykferenc/eecoin/internal/peer/domain/peer"
)

type SavePeersCommand struct {
	PathToFile string
}

type SavePeersCommandHandler interface {
	Handle(cmd SavePeersCommand) error
}

type savePeersCommandHandler struct {
	peersCtx peer.PeerContext
}

func NewSavePeersCommandHandler(peersCtx peer.PeerContext) SavePeersCommandHandler {
	return &savePeersCommandHandler{peersCtx: peersCtx}
}

func (h *savePeersCommandHandler) Handle(cmd SavePeersCommand) error {
	if cmd.PathToFile == "" {
		slog.Warn("Path to file is empty")
		return nil
	}

	file, err := os.Open(cmd.PathToFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return peer.PeersToFile(h.peersCtx.Peers(), file)
}
