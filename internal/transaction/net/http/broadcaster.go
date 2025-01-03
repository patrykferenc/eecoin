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

	dto := AsDTO(tx)
	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("could not marshal block: %w", err)
	}

	for _, peer := range peers {
		go func(peer string) {
			errors <- SendTransaction(body, peer)
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

func SendTransaction(body []byte, peer string) error {
	req, err := http.NewRequest(http.MethodPost, peer+transactionURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}
