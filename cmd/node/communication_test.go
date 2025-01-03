package main_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/patrykferenc/eecoin/cmd/node/nodetest"
	"github.com/patrykferenc/eecoin/internal/blockchain/inmem/persistence"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	transactionhttp "github.com/patrykferenc/eecoin/internal/transaction/net/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
TestShouldAcceptAndBroadcastMessage tests that node can accept
a message from the wallet and broadcast it to the network.

Scenario:
- Given a two-node network
- When a wallet sends a transaction to one of the nodes
- Then the node should accept the transaction and broadcast it to the network
- And the other node should receive the transaction and add it to its blockchain
*/
func TestShouldAcceptAndBroadcastMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// ctx := context.Background()
	require := require.New(t)
	assert := assert.New(t)

	// given a network
	network := nodetest.CreateNetwork(t, "10.4.0.0/16")

	// and two nodes
	ctnrAlfa := nodetest.StartContainer(t,
		nodetest.WithDefaultConfig,
		nodetest.WithPeers("http://10.4.0.3:22137"),
		nodetest.WithNoChain(),
		nodetest.WithLogging(),
		nodetest.WithNetwork(network.Name, "10.4.0.2"),
	)

	_ = nodetest.StartContainer(t,
		nodetest.WithDefaultConfig,
		nodetest.WithPeers("http://10.4.0.2:22137"),
		nodetest.WithNoChain(),
		nodetest.WithLogging(),
		nodetest.WithNetwork(network.Name, "10.4.0.3"),
	)

	// and a wallet
	hexPrivateGenesisAddrKey := "de9f33092050cd28e2b9382c0cb42c0d3f25c7a0254db909f10e5263924a35d9"
	dBytes, err := hex.DecodeString(hexPrivateGenesisAddrKey)
	require.NoError(err)

	// with parameters
	curve := elliptic.P256()
	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = curve
	privKey.D = new(big.Int).SetBytes(dBytes)

	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(dBytes)
	if privKey.PublicKey.X == nil || privKey.PublicKey.Y == nil {
		require.Fail("Failed to derive public key")
	}

	// and a public key
	publicMarshaled, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	require.NoError(err)

	// that is the genesis address
	senderAddr := hex.EncodeToString(publicMarshaled)
	require.Equal(transaction.GENESIS_ADDRESS, senderAddr)

	// and given receiver
	privateReceiver, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // only used to generate a public key
	require.NoError(err)
	receiverAddrRaw, err := x509.MarshalPKIXPublicKey(privateReceiver.Public())
	require.NoError(err)
	receiverAddr := hex.EncodeToString(receiverAddrRaw)

	// and relying on a node
	endpoint, err := ctnrAlfa.Endpoint(context.Background(), "")
	require.NoError(err)
	nodeUrl := "http://" + endpoint
	unspentRepo := transactionhttp.NewUnspentOutputsRepository(nodeUrl)

	// when
	tx, err := transaction.New(receiverAddr, senderAddr, 100, privKey, unspentRepo)
	require.NoError(err)

	// then
	assert.NotEmpty(tx.ID())

	// and when
	// TODO: add transaction to the network
	dto := transactionhttp.AsDTO(*tx)
	body, err := json.Marshal(dto)
	require.NoError(err)

	err = transactionhttp.SendTransaction(body, nodeUrl)
	require.NoError(err)
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
