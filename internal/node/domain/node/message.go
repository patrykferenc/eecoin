package node

import "fmt"

type PeersRepository interface {
	Get() ([]string, error)
}

var (
	ErrAllPeersFailed  = fmt.Errorf("all peers failed")
	ErrSomePeersFailed = fmt.Errorf("partial peers failed")
)

type MessageSender interface {
	SendMessage(peers []string, transaction *Transaction) error
}

type NoOpMessageSender struct{}

func (s *NoOpMessageSender) SendMessage(peers []string, transaction *Transaction) error {
	return nil
}
