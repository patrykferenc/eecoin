package nodetest

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"
	"testing"
	"time"

	dockernetwork "github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

// errLogConsumer is a log consumer that logs errors from the container.
// see the testcontainers.LogConsumer interface for more information.
type errLogConsumer struct{}

func (lc errLogConsumer) Accept(l testcontainers.Log) {
	slog.Info("Container: ", "error", string(l.Content))
}

var defaultNodeOptions = []NodeOption{
	WithDefaultConfig,
	WithNoPeers(),
	WithNoChain(),
	WithLogging(),
}

type NodeOption func(*testcontainers.ContainerRequest)

var WithDefaultConfig = WithConfigFile(filepath.Join("..", "..", "config", "example.yaml"))

func WithConfigFile(path string) NodeOption {
	return func(req *testcontainers.ContainerRequest) {
		req.Files = append(req.Files, testcontainers.ContainerFile{
			HostFilePath:      path,
			ContainerFilePath: "/etc/eecoin/config.yaml",
			FileMode:          0644,
		})
	}
}

// WithNoPeers is a NodeOption that will create a container with an empty peers file.
func WithNoPeers() NodeOption {
	return func(req *testcontainers.ContainerRequest) {
		req.Files = append(req.Files, testcontainers.ContainerFile{
			Reader:            strings.NewReader(""),
			ContainerFilePath: "/etc/eecoin/peers",
			FileMode:          0644,
		})
	}
}

func WithPeers(peers ...string) NodeOption {
	return func(req *testcontainers.ContainerRequest) {
		peerList := strings.Join(peers, "\n")
		req.Files = append(req.Files, testcontainers.ContainerFile{
			Reader:            strings.NewReader(peerList),
			ContainerFilePath: "/etc/eecoin/peers",
			FileMode:          0644,
		})
	}
}

// WithNoChain is a NodeOption that will create a container with an empty chain file - it will only contain the genesis block.
func WithNoChain() NodeOption {
	return func(req *testcontainers.ContainerRequest) {
		req.Files = append(req.Files, testcontainers.ContainerFile{
			Reader:            strings.NewReader(""),
			ContainerFilePath: "/etc/eecoin/chain",
			FileMode:          0644,
		})
	}
}

// WithLogging is a NodeOption that will log the container's output at the info level.
func WithLogging() NodeOption {
	return func(req *testcontainers.ContainerRequest) {
		req.LogConsumerCfg = &testcontainers.LogConsumerConfig{Consumers: []testcontainers.LogConsumer{&errLogConsumer{}}}
	}
}

func WithNetwork(networkName string, address string) NodeOption {
	return func(req *testcontainers.ContainerRequest) {
		req.Networks = []string{networkName}
		req.EnpointSettingsModifier = func(endpoint map[string]*dockernetwork.EndpointSettings) {
			endpoint[networkName].IPAMConfig = &dockernetwork.EndpointIPAMConfig{
				IPv4Address: address,
			}
		}
	}
}

// StartContainer starts a container with the eecoin node.
// The options will be applied to the container, and if
// no options are provided, the default options will be used.
func StartContainer(t *testing.T, opts ...NodeOption) testcontainers.Container {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:   filepath.Join("..", ".."),
			KeepImage: true,
		},
		WaitingFor: wait.ForLog("Listening on").WithStartupTimeout(3 * time.Second).WithPollInterval(1 * time.Second),
	}

	if len(opts) == 0 {
		opts = defaultNodeOptions
	}

	for _, opt := range opts {
		opt(&req)
	}

	nodeContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	testcontainers.CleanupContainer(t, nodeContainer)
	require.NoError(t, err)

	return nodeContainer
}

// CreateNetwork creates a new network with the given subnet.
func CreateNetwork(t *testing.T, subnet string) *testcontainers.DockerNetwork {
	t.Helper()
	ctx := context.Background()

	require.NotEmpty(t, subnet)

	gateway := strings.Split(subnet, "/")[0]
	gateway = strings.TrimSuffix(gateway, "0") + "1"

	newNetwork, err := network.New(ctx,
		network.WithDriver("bridge"),
		network.WithIPAM(&dockernetwork.IPAM{
			Config: []dockernetwork.IPAMConfig{
				{
					Subnet:  subnet,
					Gateway: gateway,
				},
			},
		}),
	)
	require.NoError(t, err)

	testcontainers.CleanupNetwork(t, newNetwork)

	return newNetwork
}
