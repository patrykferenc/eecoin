package nodetest

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type errLogConsumer struct{}

func (lc errLogConsumer) Accept(l testcontainers.Log) {
	slog.Info("Container: ", "error", string(l.Content))
}

func StartContainer(t *testing.T) testcontainers.Container {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    filepath.Join("..", ".."),
			Dockerfile: "Dockerfile",
			KeepImage:  true,
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      filepath.Join("..", "..", "config", "example.yaml"),
				ContainerFilePath: "/etc/eecoin/config.yaml",
				FileMode:          0644,
			},
			{
				Reader:            strings.NewReader(""),
				ContainerFilePath: "/etc/eecoin/chain",
				FileMode:          0644,
			},
			{
				Reader:            strings.NewReader(""),
				ContainerFilePath: "/etc/eecoin/peers",
				FileMode:          0644,
			},
		},
		WaitingFor:     wait.ForLog("Listening on"),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{Consumers: []testcontainers.LogConsumer{&errLogConsumer{}}},
	}

	nodeContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	testcontainers.CleanupContainer(t, nodeContainer)
	require.NoError(t, err)

	return nodeContainer
}
