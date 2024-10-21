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

	for _, p := range peers {
		b, err := json.Marshal(transaction)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, p+messagePath, bytes.NewReader(b))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to send message: unexpected status code: %d", resp.StatusCode)
		}
	}

	return nil
}
