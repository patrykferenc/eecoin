package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/peer/query"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

var url = "/transaction"

type Broadcaster struct {
	client   *http.Client
	getPeers query.GetPeers
}

func NewBroadcaster(client *http.Client, p query.GetPeers) *Broadcaster {
	if client == nil {
		client = http.DefaultClient // todo: tidier api
	}
	return &Broadcaster{client: client, getPeers: p}
}

func (b *Broadcaster) Broadcast(tx transaction.Transaction) error {
	peers, err := b.getPeers.Get()
	if err != nil {
		return fmt.Errorf("could not get peers: %w", err)
	}
	errors := make(chan error, len(peers))

	dto := asDTO(tx)
	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("could not marshal block: %w", err)
	}

	for _, peer := range peers {
		go func(peer string) {
			req, err := http.NewRequest(http.MethodPost, peer+url, bytes.NewReader(body))
			if err != nil {
				errors <- err
				return
			}

			_, err = b.client.Do(req)
			if err != nil {
				errors <- err
				return
			}

			errors <- nil
		}(peer)
	}

	for range peers {
		err := <-errors
		if err != nil {
			slog.Warn("could not broadcast transaction", "error", err)
		}
	}

	return nil // todo: return error if all peers failed
}

type inputDTO struct {
	OutputID    string `json:"output_id"`
	OutputIndex int    `json:"output_index"`
	Signature   string `json:"signature"`
}

type outputDTO struct {
	Amount  int    `json:"amount"`
	Address string `json:"address"`
}

type transactionDTO struct {
	ID      string      `json:"id"`
	Inputs  []inputDTO  `json:"inputs"`
	Outputs []outputDTO `json:"outputs"`
}

func asDTO(tx transaction.Transaction) transactionDTO {
	inputs := make([]inputDTO, len(tx.Inputs()))
	for i, in := range tx.Inputs() {
		inputs[i] = inputDTO{
			OutputID:    in.OutputID().String(),
			OutputIndex: in.OutputIndex(),
			Signature:   in.Signature(),
		}
	}

	outputs := make([]outputDTO, len(tx.Outputs()))
	for i, out := range tx.Outputs() {
		outputs[i] = outputDTO{
			Amount:  out.Amount(),
			Address: out.Address(),
		}
	}

	return transactionDTO{
		ID:      tx.ID().String(),
		Inputs:  inputs,
		Outputs: outputs,
	}
}
