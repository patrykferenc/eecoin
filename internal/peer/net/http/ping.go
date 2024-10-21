package http

import (
	"fmt"
	"net/http"

	"github.com/patrykferenc/eecoin/internal/peer/command"
)

const pingPath = "/ping"

type PingClient struct {
	acceptPingHandler command.AcceptPingHandler
	client            http.Client
}

func (p *PingClient) Ping(targetHost string) error {
	req, err := http.NewRequest(http.MethodGet, targetHost+pingPath, nil)
	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
