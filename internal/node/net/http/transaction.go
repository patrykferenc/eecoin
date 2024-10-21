package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

func postTransaction(acceptMessageHandler command.AcceptMessageHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transaction *node.Transaction

		err := json.NewDecoder(r.Body).Decode(&transaction)
		if err != nil {
			slog.Warn("Failed to decode transaction", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cmd, err := command.NewAcceptMessage(transaction)
		if err != nil {
			slog.Warn("Failed to create command", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := acceptMessageHandler.Handle(cmd); err != nil {
			slog.Warn("Failed to handle command", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError) // TODO: better error handling
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
