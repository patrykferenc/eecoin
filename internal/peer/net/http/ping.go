package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/patrykferenc/eecoin/internal/peer/command"
)

const pingPath = "/ping"

type PingClient struct {
	acceptPingHandler command.AcceptPingHandler
	client            http.Client
}

func NewPingClient(acceptPingHandler command.AcceptPingHandler) *PingClient {
	return &PingClient{
		acceptPingHandler: acceptPingHandler,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
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
