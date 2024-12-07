package command

import "github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"

type BroadcastBlock struct {
	Index int
}

type BroadcastBlockHandler interface {
	Handle(cmd BroadcastBlock) error
}

func NewBroadcastBlockHandler(repository BlockChainRepository, broadcaster broadcaster, peers peers) BroadcastBlockHandler {
	return &broadcastBlockHandler{
		repository:  repository,
		broadcaster: broadcaster,
		peers:       peers,
	}
}

type broadcastBlockHandler struct {
	repository  BlockChainRepository
	broadcaster broadcaster
	peers       peers
}

type broadcaster interface {
	Broadcast(blockchain.Block, []string) error
}

type peers interface {
	Get() ([]string, error)
}

func (h *broadcastBlockHandler) Handle(cmd BroadcastBlock) error {
	chain := h.repository.GetChain()
	block, err := chain.GetBlock(cmd.Index)
	if err != nil {
		return err
	}

	peers, err := h.peers.Get()
	if err != nil {
		return err
	}

	err = h.broadcaster.Broadcast(block, peers)
	if err != nil {
		return err
	}

	return nil
}
