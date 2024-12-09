package query

import "github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"

type GetTransactionPool interface {
	GetAll() []transaction.Transaction
}
