package peer

type PingSender interface {
	Ping(targetHost string) error
}

type PeerContext interface {
	Peers() *Peers
}
