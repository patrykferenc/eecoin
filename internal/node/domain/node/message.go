package node

type PeersRepository interface {
	Get() ([]string, error)
}

type MessageSender interface {
	SendMessage(peers []string, transaction *Transaction) error
}

type NoOpMessageSender struct{}

func (s *NoOpMessageSender) SendMessage(peers []string, transaction *Transaction) error {
	return nil
}
