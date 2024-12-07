package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
)

var url = "/block"

type Broadcaster struct {
	client *http.Client
}

func NewBroadcaster(client *http.Client) *Broadcaster {
	if client == nil {
		client = http.DefaultClient // todo: tidier api
	}
	return &Broadcaster{client: client}
}

func (b *Broadcaster) Broadcast(block blockchain.Block, peers []string) error {
	errors := make(chan error, len(peers))

	dto := asDTO(block)
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
			slog.Warn("could not broadcast block", "error", err)
		}
	}

	return nil // todo: return error if all peers failed
}

//	type Block struct {
//		Index          int
//		TimestampMilis int64
//		ContentHash    uint64
//		PrevHash       uint64
//		Transactions   []TransactionID
//		Challenge      Challenge
//	}
type blockDTO struct {
	Index          int      `json:"index"`
	TimestampMilis int64    `json:"timestamp"`
	ContentHash    uint64   `json:"content_hash"`
	PrevHash       uint64   `json:"prev_hash"`
	Transactions   []string `json:"transactions"` // TODO#30
}

func asDTO(block blockchain.Block) blockDTO {
	transactions := make([]string, len(block.Transactions))
	for i, transaction := range block.Transactions {
		transactions[i] = string(transaction)
	}

	return blockDTO{
		Index:          block.Index,
		TimestampMilis: block.TimestampMilis,
		ContentHash:    block.ContentHash,
		PrevHash:       block.PrevHash,
		Transactions:   transactions,
	}
}
