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
