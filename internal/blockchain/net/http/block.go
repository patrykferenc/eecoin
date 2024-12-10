package http

import (
	"encoding/json"
	t "github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
)

func postBlock(addBlockHandler command.AddBlockHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto blockDTO

		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			slog.Warn("failed to decode block JSON", "error", err)
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		transactions := make([]t.Transaction, len(dto.Transactions)) // TODO#30
		for i, tx := range dto.Transactions {
			translated, err := asModel(tx)
			if err != nil {
				slog.Warn("failed to decode transaction", "error", err)
				http.Error(w, "invalid body or faulty decoding method", http.StatusInternalServerError)
				return
			}
			transactions[i] = *translated
		}

		if err := addBlockHandler.Handle(command.AddBlock{
			ToAdd: blockchain.Block{
				Index:          dto.Index,
				TimestampMilis: dto.TimestampMilis,
				ContentHash:    dto.ContentHash,
				PrevHash:       dto.PrevHash,
				Transactions:   transactions,
			},
		}); err != nil {
			slog.Warn("failed to add block to chain", "error", err)
			http.Error(w, "failed to add block to chain", http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
	}
}
