package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/stretchr/testify/assert"
)

func TestUnspentClient_New(t *testing.T) {
	remote := "http://localhost:8080"
	client := NewUnspentOutputsRepository(remote)
	assert.NotNil(t, client)
}

func TestUnspentClient_GetAll_WhenError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	repo := NewUnspentOutputsRepository(mockServer.URL)

	// when
	_, err := repo.GetAll()

	// then
	assert.Error(err)
}

func TestUnspentClient_GetAll_WhenEmptyTransactionsWhenNoneExist(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("{\"unspentOutputs\": [], \"count\": 0}")); err != nil {
			t.Fatal(err)
		}
	}))

	repo := NewUnspentOutputsRepository(mockServer.URL)

	// when
	transactions, err := repo.GetAll()

	// then
	assert.NoError(err)
	assert.Empty(transactions)
}

func TestUnspentClient_GetAll_WhenTransactionExists(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("{\"outputs\": [{\"output_id\": \"123\", \"output_index\": 0, \"amount\": 100, \"address\": \"someAddress\"}], \"count\": 1}")); err != nil {
			t.Fatal(err)
		}
	}))

	repo := NewUnspentOutputsRepository(mockServer.URL)

	// when
	transactions, err := repo.GetAll()

	// then
	assert.NoError(err)
	assert.Len(transactions, 1)
	assert.Equal(transaction.ID("123"), transactions[0].OutputID())
	assert.Equal(0, transactions[0].OutputIndex())
	assert.Equal(100, transactions[0].Amount())
	assert.Equal("someAddress", transactions[0].Address())
}
