package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/query"
)

func Route(r chi.Router, addBlock command.AddBlockHandler, chain query.GetChain) {
	r.Post("/block", postBlock(addBlock))
	r.Get("/chain", getChain(chain))
}
