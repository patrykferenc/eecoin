package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/patrykferenc/eecoin/internal/node/command"
)

func Route(r chi.Router, acceptClient command.AcceptClientMessageHandler, accept command.AcceptMessageHandler) {
	r.Post("/client/message", postClientMessage(acceptClient))
	r.Post("/transaction", postTransaction(accept))
}
