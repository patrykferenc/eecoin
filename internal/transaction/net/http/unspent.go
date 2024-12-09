package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/transaction/query"
)

const unspentURL = "/unspent"

func getUnspent(q query.GetUnspentOutputs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oo, err := q.Get()
		if err != nil {
			slog.Warn("failed to get unspent outputs", "error", err)
			http.Error(w, "failed to get unspent outputs", http.StatusInternalServerError)
			return
		}

		// Respond with unspent outputs
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(oo)
	}
}
