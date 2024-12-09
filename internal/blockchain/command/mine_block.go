package command

import (
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	ev "github.com/patrykferenc/eecoin/internal/common/event"
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
	repository BlockChainRepository
	publisher  ev.Publisher
}

func NewMineBlockHandler(repository BlockChainRepository, publisher ev.Publisher) MineBlockHandler {
	return &mineBlockHandler{
		repository: repository,
		publisher:  publisher,
	}
}

func (h *mineBlockHandler) Handle(cmd MineBlock) {
	var chain = h.repository.GetChain()
	var previousBlock = chain.GetLast()
	var c, err = blockchain.NewChallenge(previousBlock.Challenge.Difficulty, previousBlock.Challenge.TimeCapMillis)

	// TODO actual transactions should be here
	var transactionsStub []blockchain.TransactionID
	if err != nil {
		slog.Error("Error creating challenge", "error", err)
	}
	for {
		currentTime := time.Now().UnixMilli()
		err := c.RollNonce(previousBlock, transactionsStub, currentTime)
		if err != nil {
			slog.Error("Error rolling nonce", "error", err)
			continue
		}
		if c.MatchesDifficulty() {
			b, err := chain.NewBlock(currentTime, transactionsStub, c)
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
		}
		select {
		case <-cmd.InterruptChannel:
			chain = h.repository.GetChain()
			previousBlock = chain.GetLast()
		}
	}
}
