package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingSender(t *testing.T) {
	// given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status": "healthy"}`)); err != nil {
			t.Fatal(err)
		}
	}))
	defer mockServer.Close()

	pingClient := &PingClient{}

	// when
	err := pingClient.Ping(mockServer.URL)

	// then
	assert.NoError(t, err)
}

func TestPingSender_error(t *testing.T) {
	// given
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	pingClient := &PingClient{}

	// when
	err := pingClient.Ping(mockServer.URL)

	// then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
