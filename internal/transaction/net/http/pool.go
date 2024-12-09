package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/transaction/query"
)

const transactionPoolURL = "/pool"

func getTransactionPool(q query.GetTransactionPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactions := q.GetAll()

		var dto transactionPoolDTO
		transactionsDTO := make([]transactionDTO, len(transactions))
		for i, t := range transactions {
			transactionsDTO[i] = asDTO(t)
		}
		dto.Transactions = transactionsDTO
		dto.Count = len(transactions)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dto); err != nil {
			slog.Warn("failed to encode transaction pool", "error", err)
			http.Error(w, "failed to encode transaction pool", http.StatusInternalServerError)
			return
		}
	}
}

type transactionPoolDTO struct {
	Transactions []transactionDTO `json:"transactions"`
	Count        int              `json:"count"`
}
