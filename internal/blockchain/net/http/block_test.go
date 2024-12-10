package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrykferenc/eecoin/internal/blockchain/command"
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/stretchr/testify/assert"
)

// // MockService mocks the blockchain.Service interface
type mockHandler struct {
	called   int
	err      error
	expected *blockchain.Block
}

func (m *mockHandler) Handle(command command.AddBlock) error {
	m.called++
	if m.expected != nil {
		if m.expected.Index != command.ToAdd.Index {
			return fmt.Errorf("expected %d, got %d", m.expected.Index, command.ToAdd.Index)
		}
		if m.expected.TimestampMilis != command.ToAdd.TimestampMilis {
			return fmt.Errorf("expected %d, got %d", m.expected.TimestampMilis, command.ToAdd.TimestampMilis)
		}
		if m.expected.ContentHash != command.ToAdd.ContentHash {
			return fmt.Errorf("expected %s, got %s", m.expected.ContentHash, command.ToAdd.ContentHash)
		}
		if m.expected.PrevHash != command.ToAdd.PrevHash {
			return fmt.Errorf("expected %s, got %s", m.expected.PrevHash, command.ToAdd.PrevHash)
		}
	}
	return m.err
}

func TestBlockHandler_HandleBlockPost(t *testing.T) {
	// given
	sampleTransaction := transaction.Transaction{Id: "skibidi"}
	sampleTransactionDTO := transactionDTO{ID: "skibidi"}
	mockBlockDTO := blockDTO{
		Index:          1,
		TimestampMilis: 1630000000000,
		ContentHash:    "12345",
		PrevHash:       "54321",
		Transactions:   []transactionDTO{sampleTransactionDTO},
	}
	mockBlock := blockchain.Block{
		Index:          mockBlockDTO.Index,
		TimestampMilis: mockBlockDTO.TimestampMilis,
		ContentHash:    mockBlockDTO.ContentHash,
		PrevHash:       mockBlockDTO.PrevHash,
		Transactions:   []transaction.Transaction{sampleTransaction},
	}

	t.Run("Successful Block Post", func(t *testing.T) {
		// given
		handler := &mockHandler{
			expected: &mockBlock,
		}
		body, _ := json.Marshal(mockBlockDTO)
		req := httptest.NewRequest(http.MethodPost, "/block", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// when
		postBlock(handler)(rec, req)

		// then
		assert.Equal(t, http.StatusOK, rec.Code, "Expected HTTP 200 OK")
		assert.Equal(t, 1, handler.called, "Expected handler to be called once")
	})

	t.Run("Invalid JSON Body", func(t *testing.T) {
		// given
		handler := &mockHandler{}
		req := httptest.NewRequest(http.MethodPost, "/block", bytes.NewReader([]byte("invalid json")))
		rec := httptest.NewRecorder()

		// when
		postBlock(handler)(rec, req)

		// then
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Expected HTTP 400 Bad Request")
		assert.Equal(t, 0, handler.called, "Expected handler not to be called")
	})

	t.Run("Service Error", func(t *testing.T) {
		// given
		handler := &mockHandler{
			err: fmt.Errorf("service error"),
		}
		body, _ := json.Marshal(mockBlockDTO)
		req := httptest.NewRequest(http.MethodPost, "/block", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		// when
		postBlock(handler)(rec, req)

		// then
		assert.Equal(t, http.StatusInternalServerError, rec.Code, "Expected HTTP 500 Internal Server Error")
		assert.Equal(t, 1, handler.called, "Expected handler to be called once")
	})
}
