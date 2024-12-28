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

var transactionURL = "/transaction"

type Broadcaster struct {
	client   *http.Client
	getPeers query.GetPeers
}

func NewBroadcaster(p query.GetPeers) *Broadcaster {
	return &Broadcaster{client: http.DefaultClient, getPeers: p}
}

func (b *Broadcaster) Broadcast(tx transaction.Transaction) error {
	peers, err := b.getPeers.Get()
	if err != nil {
		return fmt.Errorf("could not get peers: %w", err)
	}
	if len(peers) == 0 {
		return fmt.Errorf("no peers to broadcast to")
	}
	errors := make(chan error, len(peers))

	dto := asDTO(tx)
	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("could not marshal block: %w", err)
	}

	for _, peer := range peers {
		go func(peer string) {
			req, err := http.NewRequest(http.MethodPost, peer+transactionURL, bytes.NewReader(body))
			if err != nil {
				errors <- err
				return
			}

			res, err := b.client.Do(req)
			if err != nil {
				errors <- err
				return
			}
			if res.StatusCode != http.StatusOK {
				errors <- fmt.Errorf("unexpected status code: %d", res.StatusCode)
				return
			}

			errors <- nil
		}(peer)
	}

	failed := 0
	for range peers {
		err := <-errors
		if err != nil {
			slog.Warn("could not broadcast transaction", "error", err)
			failed++
		}
	}

	if failed == len(peers) {
		return fmt.Errorf("could not broadcast transaction to any peer")
	}

	return nil
}
