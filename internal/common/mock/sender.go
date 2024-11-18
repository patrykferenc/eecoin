package mock

import "github.com/patrykferenc/eecoin/internal/node/domain/node"

type MessageSender struct {
	Err error
}

func (m MessageSender) SendMessage(peers []string, transaction *node.Transaction) error {
	return m.Err
}
