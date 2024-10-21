package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/patrykferenc/eecoin/internal/peer/command"
)

func Route(router *chi.Mux, acceptPing command.AcceptPingHandler) {
	router.Get("/ping", getPing(acceptPing))
}
