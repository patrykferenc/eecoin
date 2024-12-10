package command

import "github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"

type BroadcastBlock struct {
	Block blockchain.Block
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
	_ = h.repository.GetChain() // TODO for lint xd
	peers, err := h.peers.Get()
	if err != nil {
		return err
	}

	err = h.broadcaster.Broadcast(cmd.Block, peers)
	if err != nil {
		return err
	}

	return nil
}
