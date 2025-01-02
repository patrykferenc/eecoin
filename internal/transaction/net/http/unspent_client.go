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
	remote string
}

func (u *unspentOutputsRepository) GetByOutputIDAndIndex(outputID transaction.ID, outputIndex int) (transaction.UnspentOutput, error) {
	allOutputs, err := u.GetAll()
	if err != nil {
		return transaction.UnspentOutput{}, fmt.Errorf("failed to get unspent outputs: %w", err)
	}

	for _, output := range allOutputs {
		if output.OutputID() == outputID && output.OutputIndex() == outputIndex {
			return output, nil
		}
	}

	return transaction.UnspentOutput{}, fmt.Errorf("output with ID %s and index %d not found on remote", outputID, outputIndex)
}

func NewUnspentOutputsRepository(remote string) *unspentOutputsRepository {
	return &unspentOutputsRepository{remote: remote}
}

func (u *unspentOutputsRepository) GetAll() ([]transaction.UnspentOutput, error) {
	var dto query.UnspentOutputs

	resp, err := u.client.Get(u.remote + unspentURL)
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
	return fmt.Errorf("not implemented for %v", unspentOutputs) // TODO#30
}

func (u *unspentOutputsRepository) Get(peers []string) ([]transaction.UnspentOutput, error) { // TODO#30 - possibly tidy it up
	var allErrors []error

	for _, peer := range peers {
		resp, err := u.client.Get(peer + unspentURL)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to get unspent outputs from %s: %w", peer, err))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			allErrors = append(allErrors, fmt.Errorf("failed to get unspent outputs from %s: %v", peer, resp.Status))
			resp.Body.Close()
			continue
		}

		defer resp.Body.Close()

		var dto query.UnspentOutputs
		if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to decode unspent outputs from %s: %w", peer, err))
			continue
		}

		return dto.ToModel(), nil
	}
	return nil, fmt.Errorf("failed to fetch unspent outputs from all peers: %v", allErrors)
}
