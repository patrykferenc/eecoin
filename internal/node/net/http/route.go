package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/query"
)

func Route(r chi.Router, acceptClient command.AcceptClientMessageHandler, accept command.AcceptMessageHandler, chain query.GetChain) {
	r.Post("/v0/client/message", postClientMessage(acceptClient))
	r.Post("/v0/transaction", postTransaction(accept))
	r.Get("/chain", getChain(chain))
}
