package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/patrykferenc/eecoin/internal/transaction/command"
	"github.com/patrykferenc/eecoin/internal/transaction/query"
)

func Route(r chi.Router, addTransaction command.AddTransactionHandler, unspent query.GetUnspentOutputs) {
	r.Post(transactionURL, postTransaction(addTransaction))
	r.Get(unspentURL, getUnspent(unspent))
}
