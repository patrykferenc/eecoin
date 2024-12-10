package main

import (
	"crypto/x509"
	"github.com/patrykferenc/eecoin/internal/blockchain/inmem/persistence"
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
	"log/slog"
	"net/http"
	"os"
	"time"

	bc "github.com/patrykferenc/eecoin/internal/blockchain"
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	blockchaincommand "github.com/patrykferenc/eecoin/internal/blockchain/command"
	blockchainHttp "github.com/patrykferenc/eecoin/internal/blockchain/net/http"
	"github.com/patrykferenc/eecoin/internal/common/config"
	"github.com/patrykferenc/eecoin/internal/common/event"
	peercntr "github.com/patrykferenc/eecoin/internal/peer"
	peercommand "github.com/patrykferenc/eecoin/internal/peer/command"
	peerhttp "github.com/patrykferenc/eecoin/internal/peer/net/http"
	transactionhttp "github.com/patrykferenc/eecoin/internal/transaction/net/http"
)

func main() {
	slog.Info("Starting Eecoin node")

	cfg, err := readConfig()
	if err != nil {
		slog.Error("Failed to read config", "error", err)
		return
	}
	if cfg.Persistence.SelfKey == "" {
		key, err := wallet.NewEcdsaKey()
		if err != nil {
			slog.Error("Failed to generate key", "error", err)
			return
		}
		marshalled, err := x509.MarshalPKIXPublicKey(key.Public)
		if err != nil {
			slog.Error("Failed to marshal public key", "error", err)
			return
		}
		cfg.Persistence.SelfKey = string(marshalled)
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
	go schedulePersistChain(cfg, container.blockChainComponent.Queries.GetChain.Get())
	go schedulePing(cfg, container.peerComponent)
	go scheduleMining(container.blockChainComponent, container.interruptionChanel)

	go pubSub(container)

	go sync(container)

	if err := listenAndServe(container); err != nil {
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

func listenAndServe(container *Container) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	peerhttp.Route(r, container.peerComponent.Commands.AcceptPing)
	blockchainHttp.Route(r, container.blockChainComponent.Commands.AddBlock, container.blockChainComponent.Queries.GetChain)
	transactionhttp.Route(
		r,
		container.transactionComponent.Commands.AddTransactionHandler,
		container.transactionComponent.Queries.GetUnspentOutputs,
		container.transactionComponent.Queries.GetTransactionPool,
	)

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

func scheduleMining(blockchainComponent *bc.Component, interrupt chan bool) {
	h := blockchainComponent.Commands.MineBlock
	for {
		h.Handle(blockchaincommand.MineBlock{InterruptChannel: interrupt})
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
		err := persistence.Persist(chain, cfg.Persistence.ChainFilePath)
		if err != nil {
			slog.Error("Failed to persist blockchain", "error", err)
		}
	}
}

func pubSub(cntr *Container) {
	handlers := map[string]func(event.Event) error{
		"x.block.added": func(e event.Event) error {
			data, ok := e.Data().(blockchain.NewBlockAddedEvent)
			if !ok {
				slog.Error("Invalid event data")
				return nil
			}
			cmd := blockchaincommand.AddBlock{ToAdd: data.Block}
			err := cntr.blockChainComponent.Commands.AddBlock.Handle(cmd)
			if err != nil {
				slog.Error("Failed to handle AddBlock command", "error", err)
			}

			//c <- true
			return nil
		},
	}

	cntr.broker.RouteAll(handlers)
}

func sync(cntr *Container) {
	time.Sleep(5 * time.Second)
	err := cntr.transactionComponent.Application.TransactionUpdater.UpdateFromRemote()
	if err != nil {
		slog.Error("Failed to sync", "error", err)
	} else {
		slog.Info("Synced")
	}
}
