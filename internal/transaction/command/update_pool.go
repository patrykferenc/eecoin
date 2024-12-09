package command

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type UpdatePool struct{}

// UpdatePoolHandler updates the pool with the latest unspent outputs
type UpdatePoolHandler interface {
	Handle(cmd UpdatePool) error
}

type updatePoolHandler struct {
	pool transaction.Pool
	repo transaction.UnspentOutputRepository
}

func NewUpdatePoolHandler(pool transaction.Pool, repo transaction.UnspentOutputRepository) UpdatePoolHandler {
	return &updatePoolHandler{
		pool: pool,
		repo: repo,
	}
}

func (h *updatePoolHandler) Handle(cmd UpdatePool) error {
	unspent, err := h.repo.GetAll()
	if err != nil {
		return fmt.Errorf("could not get unspent outputs: %w", err)
	}

	return h.pool.Update(unspent)
}
