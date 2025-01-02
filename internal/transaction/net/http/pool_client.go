package http

import (
	"encoding/json"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type TransactionPoolClient struct {
	client http.Client
}

func (t *TransactionPoolClient) Get(peers []string) ([]transaction.Transaction, error) {
	var errs []error

	for _, peer := range peers {
		resp, err := t.client.Get(peer + poolURL)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			errs = append(errs, err)
			continue
		}
		defer resp.Body.Close()

		var dto transactionPoolDTO
		if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
			errs = append(errs, err)
			continue
		}

		var transactions []transaction.Transaction
		for _, t := range dto.Transactions {
			tx, err := asModel(t)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			transactions = append(transactions, *tx)
		}

		return transactions, nil
	}

	return nil, errs[0] // TODO: abomination
}
