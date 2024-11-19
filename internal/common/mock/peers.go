package mock

type Peers struct {
	Peers []string
	Err   error
}

func NewPeersFailing(err error) Peers {
	return Peers{Err: err}
}

func NewPeers(peers []string) Peers {
	return Peers{Peers: peers}
}

func (m Peers) Get() ([]string, error) {
	return m.Peers, m.Err
}
