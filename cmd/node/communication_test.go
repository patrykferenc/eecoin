package main_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/patrykferenc/eecoin/cmd/node/nodetest"
	"github.com/patrykferenc/eecoin/internal/blockchain/inmem/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
TestShouldAcceptAndBroadcastMessage tests that node can accept
a message from the wallet and broadcast it to the network.

Scenario:
- Given a node
...
*/
func TestShouldAcceptAndBroadcastMessage(t *testing.T) {
	// TODO: Implement E2E test
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	t.Skip("TODO: Implement E2E test")
}

/*
TestShouldStartFresh tests that the node starts fresh, meaning that it has an empty blockchain (only genesis block).

Scenario:
- Given a node
- When starting it (defaults have empty peers and empty blockchain file)
- Then it should contain only the genesis block
*/
func TestShouldStartFresh(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	ctx := context.Background()
	require := require.New(t)
	assert := assert.New(t)

	// given container, when building the project, then it should succeed
	ctnr := nodetest.StartContainer(t)
	endpoint, err := ctnr.Endpoint(ctx, "")
	require.NoError(err)
	url := "http://" + endpoint

	// and blockchain should contain genesis block
	req, err := http.NewRequest("GET", url+"/chain", nil)
	require.NoError(err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(err)

	// then it should return 200 OK
	assert.Equal(http.StatusOK, res.StatusCode)
	// and it should contain genesis block
	defer res.Body.Close()
	var dtoChain persistence.ChainDto
	err = json.NewDecoder(res.Body).Decode(&dtoChain)
	require.NoError(err)
	require.Len(dtoChain.Blocks, 1)
	// and
	genesisBlock := dtoChain.Blocks[0]
	assert.EqualValues(0, genesisBlock.Index)
	assert.EqualValues(time.Date(2024, 11, 16, 20, 23, 0, 0, time.UTC).UnixMilli(), genesisBlock.TimestampMilis)
	assert.EqualValues("", genesisBlock.PrevHash)
}
