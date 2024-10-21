package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrykferenc/eecoin/internal/peer"
	"github.com/patrykferenc/eecoin/internal/peer/command"
	peerhttp "github.com/patrykferenc/eecoin/internal/peer/net/http"
)

func main() {
	slog.Info("Starting Eecoin node...")

	container, err := newContainer()
	if err != nil {
		slog.Error("Failed to create container", "error", err)
		return
	}
	slog.Info("Context started...")

	if os.Getenv("EECOIN_SAVE_PEERS") != "" {
		slog.Info("Will be saving peers to file")
		go scheduleSave(container.peerComponent)
	}

	go schedulePing(container.peerComponent)

	if err := listenAndServe(container.peerComponent); err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		return
	}
}

func listenAndServe(peerComponent *peer.Component) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	peerhttp.Route(r, peerComponent.Commands.AcceptPing)

	return http.ListenAndServe(":22137", r)
}

func schedulePing(peerComponent *peer.Component) {
	handler := peerComponent.Commands.SendPing

	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		handler.Handle(command.SendPingCommand{})
	}
}

func scheduleSave(peerComponent *peer.Component) {
	handler := peerComponent.Commands.SavePeers

	ticker := time.NewTicker(10 * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		err := handler.Handle(command.SavePeersCommand{PathToFile: "/etc/eecoin/peers"})
		if err != nil {
			slog.Error("Failed to save peers", "error", err)
		}
	}
}
