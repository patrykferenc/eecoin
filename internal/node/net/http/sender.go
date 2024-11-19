package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

const messagePath = "/transaction"

type Sender struct {
	client *http.Client
}

func NewSender() *Sender {
	return &Sender{
		client: &http.Client{},
	}
}

// TODO: retries, failure handling etc
func (s *Sender) SendMessage(peers []string, transaction *node.Transaction) error {
	if transaction == nil {
		return fmt.Errorf("transaction is required")
	}
	if len(peers) == 0 {
		slog.Warn("No peers to send message to, will skip without error")
		return nil
	}

	b, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}
	failed := 0
	for _, p := range peers {
		req, err := http.NewRequest(http.MethodPost, p+messagePath, bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err) // critical issue
		}

		resp, err := s.client.Do(req)
		if err != nil {
			slog.Warn("failed to send message to %s: %v", p, err)
			failed++
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			slog.Warn("failed to send message to %s: status code: %d", p, resp.StatusCode)
			failed++
		}
	}

	if failed > 0 {
		if failed == len(peers) {
			return node.ErrAllPeersFailed
		}

		slog.Warn("Some peers failed to receive message", "count", failed)
		return node.ErrSomePeersFailed
	}

	return nil
}
