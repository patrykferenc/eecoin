package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/patrykferenc/eecoin/internal/transaction/query"
)

type unspentOutputsRepository struct {
	client http.Client
}

func (u *unspentOutputsRepository) GetAll() ([]transaction.UnspentOutput, error) {
	var dto query.UnspentOutputs

	resp, err := u.client.Get("/unspent")
	if err != nil {
		return nil, fmt.Errorf("failed to get unspent outputs: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get unspent outputs: %v", resp.Status)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("failed to decode unspent outputs: %w", err)
	}

	return dto.ToModel(), nil
}

func (u *unspentOutputsRepository) GetByAddress(address string) ([]transaction.UnspentOutput, error) {
	outputs, err := u.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get unspent outputs: %w", err)
	}

	var uos []transaction.UnspentOutput
	for _, output := range outputs {
		if output.Address() == address {
			uos = append(uos, output)
		}
	}

	return uos, nil
}

func (u *unspentOutputsRepository) Set(unspentOutputs []transaction.UnspentOutput) error {
	return fmt.Errorf("not implemented for %v", unspentOutputs)
}
