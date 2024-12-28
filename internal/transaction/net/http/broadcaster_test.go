package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrykferenc/eecoin/internal/common/mock"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction/transactiontest"
	"github.com/stretchr/testify/assert"
)

func TestBroadcastSuccessful(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/transaction" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()
	peers := []string{mockServer.URL}
	broadcaster := NewBroadcaster(mock.NewPeers(peers))
	// and
	tx, err := transactiontest.NewTransaction()
	assert.NoError(err)

	// when
	err = broadcaster.Broadcast(*tx)

	// then
	assert.NoError(err)
}

func TestBroadcastAllPeersFailed(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()
	peers := []string{mockServer.URL}
	broadcaster := NewBroadcaster(mock.NewPeers(peers))
	// and
	tx, err := transactiontest.NewTransaction()
	assert.NoError(err)

	// when
	err = broadcaster.Broadcast(*tx)

	// then
	assert.Error(err)
}

func TestBroadcastOnePeerFailed(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// given
	mockServer1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer1.Close()
	mockServer2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer2.Close()
	peers := []string{mockServer1.URL, mockServer2.URL}
	broadcaster := NewBroadcaster(mock.NewPeers(peers))
	// and
	tx, err := transactiontest.NewTransaction()
	assert.NoError(err)

	// when
	err = broadcaster.Broadcast(*tx)

	// then
	assert.NoError(err) // should not fail if at least one peer succeeded
}

func TestBroadcastNoPeers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// given
	broadcaster := NewBroadcaster(mock.NewPeers([]string{}))
	// and
	tx, err := transactiontest.NewTransaction()
	assert.NoError(err)

	// when
	err = broadcaster.Broadcast(*tx)

	// then
	assert.Error(err)
}
