package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
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
	Index          int              `json:"index"`
	TimestampMilis int64            `json:"timestamp"`
	ContentHash    string           `json:"content_hash"`
	PrevHash       string           `json:"prev_hash"`
	Transactions   []transactionDTO `json:"transactions"` // TODO#30
	Challenge      challengeDTO     `json:"challenge"`
}

func asDTO(block blockchain.Block) blockDTO {
	transactions := make([]transactionDTO, len(block.Transactions))
	for i, trscnion := range block.Transactions {
		transactions[i] = transDTO(trscnion)
	}

	return blockDTO{
		Index:          block.Index,
		TimestampMilis: block.TimestampMilis,
		ContentHash:    block.ContentHash,
		PrevHash:       block.PrevHash,
		Transactions:   transactions,
		Challenge:      challengeModelToDTO(block.Challenge),
	}
}

type inputDTO struct {
	OutputID    string `json:"output_id"`
	OutputIndex int    `json:"output_index"`
	Signature   string `json:"signature"`
}

func (i inputDTO) asInput() *transaction.Input {
	o := transaction.NewInput(transaction.ID(i.OutputID), i.OutputIndex, i.Signature)
	return &o
}

type outputDTO struct {
	Amount  int    `json:"amount"`
	Address string `json:"address"`
}

func (o outputDTO) asOutput() *transaction.Output {
	return transaction.NewOutput(o.Amount, o.Address)
}

type transactionDTO struct {
	ID      string      `json:"id"`
	Inputs  []inputDTO  `json:"inputs"`
	Outputs []outputDTO `json:"outputs"`
}

func transDTO(tx transaction.Transaction) transactionDTO {
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

func asModel(dto transactionDTO) (*transaction.Transaction, error) {
	inputs := make([]*transaction.Input, len(dto.Inputs))
	for i, in := range dto.Inputs {
		inputs[i] = in.asInput()
	}

	outputs := make([]*transaction.Output, len(dto.Outputs))
	for i, out := range dto.Outputs {
		outputs[i] = out.asOutput()
	}

	return transaction.NewFrom(inputs, outputs)
}

type challengeDTO struct {
	Difficulty    int    `json:"difficulty"`
	Nonce         uint32 `json:"nonce"`
	HashValue     string `json:"hash_value"`
	TimeCapMillis int64  `json:"time_cap_millis"`
}

func challengeModelToDTO(challenge blockchain.Challenge) challengeDTO {
	return challengeDTO{
		Difficulty:    challenge.Difficulty,
		Nonce:         challenge.Nonce,
		HashValue:     challenge.HashValue,
		TimeCapMillis: challenge.TimeCapMillis,
	}
}
