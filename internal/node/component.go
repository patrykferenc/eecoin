package node

import (
	"github.com/patrykferenc/eecoin/internal/common/event"
	"github.com/patrykferenc/eecoin/internal/node/command"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
)

type Component struct {
	Commands Commands
}

type Commands struct {
	sendMessage         command.SendMessageHandler
	acceptClientMessage command.AcceptClientMessageHandler
	acceptMessage       command.AcceptMessageHandler
	persistMessage      command.PersistMessageHandler
}

func NewComponent(publisher event.Publisher) (Component, error) {
	repo := node.NewSimpleInFlightTransactionRepository()
	seen := node.NewSimpleSeenTransactionRepository() // TODO: refactor when adding a real blockchain impl
	noOpSender := node.NoOpMessageSender{}            // TODO: inject real sender
	sendMessage, err := command.NewSendMessageHandler(repo, seen, &noOpSender, nil, publisher)
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
			sendMessage:         sendMessage,
			acceptClientMessage: acceptClientMessage,
			acceptMessage:       acceptMessage,
			persistMessage:      persistMessage,
		},
	}, nil
}
