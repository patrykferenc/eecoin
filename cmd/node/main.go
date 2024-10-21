package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	nodecntr "github.com/patrykferenc/eecoin/internal/node"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	nodehttp "github.com/patrykferenc/eecoin/internal/node/net/http"
	peercntr "github.com/patrykferenc/eecoin/internal/peer"
	peercommand "github.com/patrykferenc/eecoin/internal/peer/command"
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

	go pubSub(container)

	if err := listenAndServe(container.peerComponent, container.nodeComponent); err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		return
	}
}

func listenAndServe(peerComponent *peercntr.Component, nodeComponent *nodecntr.Component) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	peerhttp.Route(r, peerComponent.Commands.AcceptPing)
	nodehttp.Route(r, nodeComponent.Commands.AcceptClientMessage, nodeComponent.Commands.AcceptMessage)

	return http.ListenAndServe(":22137", r)
}

func schedulePing(peerComponent *peercntr.Component) {
	handler := peerComponent.Commands.SendPing

	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		handler.Handle(peercommand.SendPingCommand{})
	}
}

func scheduleSave(peerComponent *peercntr.Component) {
	handler := peerComponent.Commands.SavePeers

	ticker := time.NewTicker(10 * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		err := handler.Handle(peercommand.SavePeersCommand{PathToFile: "/etc/eecoin/peers"})
		if err != nil {
			slog.Error("Failed to save peers", "error", err)
		}
	}
}

func pubSub(cntr *container) {
	// sentMsgs := cntr.broker.Subscribe("x.message.sent") // TODO: We will implement it later, to discard msgs
	sendMsgs := cntr.broker.Subscribe("x.message.send")

	for sendMsg := range sendMsgs {
		msg, ok := sendMsg.(node.SendMessageEvent)
		if !ok {
			slog.Error("Failed to cast message to SendMessageEvent")
			continue
		}
		cmd, err := command.NewSendMessage(msg.TransactionID)
		if err != nil {
			slog.Error("Failed to create SendMessage command", "error", err)
			continue
		}
		err = cntr.nodeComponent.Commands.SendMessage.Handle(cmd)
		if err != nil {
			slog.Error("Failed to handle SendMessage command", "error", err)
		}
	}

	slog.Info("PubSub finished")
}
