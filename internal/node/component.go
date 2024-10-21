package node

import (
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/patrykferenc/eecoin/internal/node/net/http"
)

type Component struct {
	Commands Commands
}

type Commands struct {
	SendMessage         command.SendMessageHandler
	AcceptClientMessage command.AcceptClientMessageHandler
	AcceptMessage       command.AcceptMessageHandler
	PersistMessage      command.PersistMessageHandler
}

func NewComponent(publisher event.Publisher, peersRepo node.PeersRepository) (Component, error) {
	repo := node.NewSimpleInFlightTransactionRepository()
	seen := node.NewSimpleSeenTransactionRepository() // TODO: refactor when adding a real blockchain impl
	sender := http.NewSender()
	sendMessage, err := command.NewSendMessageHandler(repo, seen, sender, peersRepo, publisher)
	if err != nil {
		return Component{}, err
	}

	acceptClientMessage, err := command.NewAcceptClientMessageHandler(repo, seen, publisher)
	if err != nil {
		return Component{}, err
	}

	acceptMessage, err := command.NewAcceptMessageHandler(repo, publisher)
	if err != nil {
		return Component{}, err
	}

	persistMessage, err := command.NewPersistMessageHandler(repo, seen)
	if err != nil {
		return Component{}, err
	}

	return Component{
		Commands: Commands{
			SendMessage:         sendMessage,
			AcceptClientMessage: acceptClientMessage,
			AcceptMessage:       acceptMessage,
			PersistMessage:      persistMessage,
		},
	}, nil
}
