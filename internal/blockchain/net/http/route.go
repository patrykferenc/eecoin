package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/patrykferenc/eecoin/internal/blockchain/command"
)

func Route(r chi.Router, addBlock command.AddBlockHandler) {
	r.Post("/block", postBlock(addBlock))
}
