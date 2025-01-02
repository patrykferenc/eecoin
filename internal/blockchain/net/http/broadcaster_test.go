package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBroadcaster_Broadcast(t *testing.T) {
	// given setup
	var mu sync.Mutex
	calledPeers := []string{}
	receivedBodies := []blockDTO{}

	// and given peers respond with 200 OK
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/block") {
			t.Errorf("unexpected request method or URL: %s %s", r.Method, r.URL.Path)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
			http.Error(w, "bad request", http.StatusInternalServerError)
			return
		}

		var receivedBlock blockDTO
		err = json.Unmarshal(bodyBytes, &receivedBlock)
		if err != nil {
			t.Errorf("failed to decode JSON: %v", err)
			http.Error(w, "bad request", http.StatusInternalServerError)
			return
		}

		calledPeers = append(calledPeers, r.Host)
		receivedBodies = append(receivedBodies, receivedBlock)

		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// and given
	peerURLs := []string{mockServer.URL}
	broadcaster := NewBroadcaster()
	sampleTransaction, err := transaction.NewFrom(nil, nil)
	require.NoError(t, err, "NewFrom should not return an error")
	sampleChallange := blockchain.Challenge{
		Nonce:         123,
		HashValue:     "12345",
		Difficulty:    2,
		TimeCapMillis: 2,
	}
	sampleChallangeDto := challengeDTO{
		Nonce:         123,
		HashValue:     "12345",
		Difficulty:    2,
		TimeCapMillis: 2,
	}
	mockBlock := blockchain.Block{
		Index:          1,
		TimestampMilis: 123456789,
		ContentHash:    "12345",
		PrevHash:       "54321",
		Transactions:   []transaction.Transaction{*sampleTransaction},
		Challenge:      sampleChallange,
	}
	expectedTransactionDTO := transDTO(*sampleTransaction)
	expectedBody := blockDTO{
		Index:          mockBlock.Index,
		TimestampMilis: mockBlock.TimestampMilis,
		ContentHash:    mockBlock.ContentHash,
		PrevHash:       mockBlock.PrevHash,
		Transactions:   []transactionDTO{expectedTransactionDTO},
		Challenge:      sampleChallangeDto,
	}
	_ = expectedBody // TODO#30 - fix issue with marshalling ID

	// when
	err = broadcaster.Broadcast(mockBlock, peerURLs)

	// then
	assert.NoError(t, err, "Broadcast should not return an error")

	mu.Lock()
	defer mu.Unlock()

	// and then
	assert.Equal(t, len(peerURLs), len(calledPeers), "All peers should be called")
	assert.Contains(t, calledPeers, strings.TrimPrefix(mockServer.URL, "http://"), "The mock server should be called")
	assert.Equal(t, 1, len(receivedBodies), "Exactly one block body should be received")
	// assert.Equal(t, expectedBody, receivedBodies[0], "The request body should match the expected blockDTO")  // TODO#30 fix issue with marshalling ID
}
