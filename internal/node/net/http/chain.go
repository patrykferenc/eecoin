package http

import (
	"encoding/json"
	"github.com/patrykferenc/eecoin/internal/blockchain/inmem/persistence"
	"github.com/patrykferenc/eecoin/internal/node/query"
	"log/slog"
	"net/http"
)

func getChain(getChainQueryHandler query.GetChain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dtoChain := persistence.MapToDto(getChainQueryHandler.Get())
		err := json.NewEncoder(w).Encode(dtoChain)
		if err != nil {
			slog.Error("Cannot encode blockchain")
			http.NotFound(w, r)
			return
		}
	}
}
