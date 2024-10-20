package peer

import (
	"fmt"
	"io"
	"strings"
)

func PeersFromFile(file io.Reader) (*Peers, error) {
	buff := make([]byte, 1024)
	n, err := file.Read(buff)
	if err != nil {
		if err == io.EOF {
			return NewPeers(nil), nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if n == 0 {
		return NewPeers(nil), nil
	}

	hosts := strings.Split(string(buff[:n]), "\n")
	peers := make([]*Peer, len(hosts))
	for i, host := range hosts {
		peers[i] = &Peer{Host: host, Status: StatusUnknown}
	}

	return NewPeers(peers), nil
}

func PeersToFile(peers *Peers, file io.Writer) error {
	for _, peer := range peers.All() {
		_, err := fmt.Fprintln(file, peer.Host)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	return nil
}
