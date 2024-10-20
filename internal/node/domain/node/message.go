package node

type Peer string

type Peers []Peer

type PeersRepository interface {
	Get() (Peers, error)
}

type MessageSender interface {
	SendMessage(peers Peers, transaction *Transaction) error
}
