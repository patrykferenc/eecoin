package main

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrykferenc/eecoin/internal/common/config"
	"github.com/patrykferenc/eecoin/internal/common/event"
	nodecntr "github.com/patrykferenc/eecoin/internal/node"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	nodehttp "github.com/patrykferenc/eecoin/internal/node/net/http"
	peercntr "github.com/patrykferenc/eecoin/internal/peer"
	peercommand "github.com/patrykferenc/eecoin/internal/peer/command"
	peerhttp "github.com/patrykferenc/eecoin/internal/peer/net/http"
)

func main() {
	slog.Info("Starting Eecoin node")

	cfg, err := readConfig()
	if err != nil {
		slog.Error("Failed to read config", "error", err)
		return
	}

	if err := setLoggerLevel(cfg); err != nil {
		slog.Error("Failed to set logger level", "error", err)
		return
	}

	container, err := NewContainer(cfg)
	if err != nil {
		slog.Error("Failed to create container", "error", err)
		return
	}
	slog.Info("Context constructed")

	go scheduleSave(cfg, container.peerComponent)
	go schedulePersistChain(cfg, container.nodeComponent.Queries.GetChain.Get())
	go schedulePing(cfg, container.peerComponent)

	go pubSub(container)

	if err := listenAndServe(container.peerComponent, container.nodeComponent); err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
		return
	}
}

func readConfig() (*config.Config, error) {
	configPath := "/etc/eecoin/config.yaml"
	if os.Getenv("EECOIN_CONFIG") != "" {
		configPath = os.Getenv("EECOIN_CONFIG")
	}

	return config.Read(configPath)
}

func setLoggerLevel(cfg *config.Config) error {
	level, err := cfg.Log.LevelIfSet()
	if err != nil {
		return err
	}
	slog.SetLogLoggerLevel(level)
	return nil
}

func listenAndServe(peerComponent *peercntr.Component, nodeComponent *nodecntr.Component) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	peerhttp.Route(r, peerComponent.Commands.AcceptPing)
	nodehttp.Route(r, nodeComponent.Commands.AcceptClientMessage, nodeComponent.Commands.AcceptMessage, nodeComponent.Queries.GetChain)

	slog.Info("Listening on :22137")
	return http.ListenAndServe(":22137", r)
}

func schedulePing(cfg *config.Config, peerComponent *peercntr.Component) {
	handler := peerComponent.Commands.SendPing
	ticker := time.NewTicker(cfg.Peers.PingDuration)

	defer ticker.Stop()
	for range ticker.C {
		handler.Handle(peercommand.SendPingCommand{})
	}
}

func scheduleSave(cfg *config.Config, peerComponent *peercntr.Component) {
	if cfg.Peers.UpdateFileDuration == 0 {
		return
	}

	handler := peerComponent.Commands.SavePeers
	ticker := time.NewTicker(cfg.Peers.UpdateFileDuration)

	defer ticker.Stop()
	for range ticker.C {
		err := handler.Handle(peercommand.SavePeersCommand{PathToFile: cfg.Peers.FilePath})
		if err != nil {
			slog.Error("Failed to save peers", "error", err)
		}
	}
}

func schedulePersistChain(cfg *config.Config, chain blockchain.BlockChain) {
	if cfg.Persistence.UpdateFileDuration == 0 {
		return
	}

	ticker := time.NewTicker(cfg.Persistence.UpdateFileDuration)

	defer ticker.Stop()
	for range ticker.C {
		err := blockchain.Persist(chain, cfg.Persistence.ChainFilePath)
		if err != nil {
			slog.Error("Failed to persist blockchain", "error", err)
		}
	}
}

func pubSub(cntr *Container) {
	handlers := map[string]func(event.Event) error{
		"x.message.send": func(e event.Event) error {
			data, ok := e.Data().(node.SendMessageEvent)
			if !ok {
				slog.Error("Invalid event data")
				return nil
			}
			cmd, err := command.NewSendMessage(data.TransactionID)
			if err != nil {
				slog.Error("Failed to create SendMessage command", "error", err)
				return nil
			}
			err = cntr.nodeComponent.Commands.SendMessage.Handle(cmd)
			if err != nil {
				slog.Error("Failed to handle SendMessage command", "error", err)
			}
			return nil
		},
		"x.message.sent": func(e event.Event) error {
			data, ok := e.Data().(node.MessageSentEvent)
			if !ok {
				slog.Error("Invalid event data")
				return nil
			}
			cmd, err := command.NewPersistMessage(data.TransactionID)
			if err != nil {
				slog.Error("Failed to create PersistMessage command", "error", err)
				return nil
			}
			err = cntr.nodeComponent.Commands.PersistMessage.Handle(cmd)
			if err != nil {
				slog.Error("Failed to handle PersistMessage command", "error", err)
			}
			return nil
		},
	}

	cntr.broker.RouteAll(handlers)
}
