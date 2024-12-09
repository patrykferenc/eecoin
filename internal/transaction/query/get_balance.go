package query

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
)

type GetBalance interface {
	GetBalance(request GetBalanceRequest) (Balance, error)
}

type GetBalanceRequest struct {
	Address string `json:"address"`
}

type Balance struct {
	ECTS int `json:"ects"`
}

type getBalance struct {
	repo transaction.UnspentOutputRepository
}

func NewGetBalance(repo transaction.UnspentOutputRepository) GetBalance {
	return &getBalance{repo: repo}
}

func (g *getBalance) GetBalance(request GetBalanceRequest) (Balance, error) {
	if request.Address == "" {
		return Balance{}, fmt.Errorf("address is required")
	}

	outputs, err := g.repo.GetByAddress(request.Address)
	if err != nil {
		return Balance{}, err
	}
	amount := 0
	for _, o := range outputs {
		amount += o.Amount()
	}

	return Balance{ECTS: amount}, nil
}
