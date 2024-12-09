package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/transaction/command"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

const transactionURL = "/transaction"

func postTransaction(addTransactionHandler command.AddTransactionHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto transactionDTO

		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			slog.Warn("failed to decode transaction JSON", "error", err)
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		inputs := make([]*transaction.Input, len(dto.Inputs))
		for i, in := range dto.Inputs {
			inputs[i] = in.asInput()
		}
		outputs := make([]*transaction.Output, len(dto.Outputs))
		for i, out := range dto.Outputs {
			outputs[i] = out.asOutput()
		}

		if err := addTransactionHandler.Handle(command.AddTransaction{
			ProvidedID: dto.ID,
			Inputs:     inputs,
			Outputs:    outputs,
		}); err != nil {
			slog.Warn("failed to add transaction to pool", "error", err)
			http.Error(w, "failed to add transaction to pool", http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
	}
}
