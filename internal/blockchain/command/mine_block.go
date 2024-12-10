package command

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	ev "github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"log/slog"
	"time"
)

type MineBlock struct {
	InterruptChannel chan bool
}

type MineBlockHandler interface {
	Handle(cmd MineBlock)
}

type mineBlockHandler struct {
	selfAddr       string
	repository     BlockChainRepository
	publisher      ev.Publisher
	poolRepository transaction.PoolRepository
}

func NewMineBlockHandler(selfAddress string, repository BlockChainRepository, publisher ev.Publisher, poolRepository transaction.PoolRepository) MineBlockHandler {
	return &mineBlockHandler{
		repository:     repository,
		publisher:      publisher,
		poolRepository: poolRepository,
		selfAddr:       selfAddress,
	}
}

func (h *mineBlockHandler) Handle(cmd MineBlock) {
	var chain = h.repository.GetChain()
	var previousBlock = chain.GetLast()
	var c, err = blockchain.NewChallenge(previousBlock.Challenge.Difficulty, previousBlock.Challenge.TimeCapMillis)

	if err != nil {
		slog.Error("Error creating challenge", "error", err)
	}
	for {
		currentTime := time.Now().UnixMilli()
		var transactions = h.poolRepository.GetAll()
		tx, err := transaction.NewCoinbase(h.selfAddr, 10)
		transactions = append([]transaction.Transaction{*tx}, transactions...)
		if err != nil {
			slog.Error("Error creating coinbase transaction", "error", err)
			continue
		}
		err = c.RollNonce(previousBlock, transactions, currentTime)
		if err != nil {
			slog.Error("Error rolling nonce", "error", err)
			continue
		}
		if c.MatchesDifficulty() {

			b, err := chain.NewBlock(currentTime, transactions, c)
			if err != nil {
				slog.Error("Error creating new block", "error", err)
				continue
			}

			err = h.repository.PutBlock(b)
			if err != nil {
				slog.Error("Error adding new block", "error", err)
				continue
			}

			event, err := ev.New(blockchain.NewBlockAddedEvent{Block: b}, "x.block.added")
			if err != nil {
				slog.Error("Error creating event", "error", err)
				continue
			}

			err = h.publisher.Publish(event)
			if err != nil {
				slog.Error("Error publishing event", "error", err)
				continue
			}

			chain = h.repository.GetChain()
			previousBlock = chain.GetLast()
		}
		if <-cmd.InterruptChannel {
			chain = h.repository.GetChain()
			previousBlock = chain.GetLast()
		}
	}
}
