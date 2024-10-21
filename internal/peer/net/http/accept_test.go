package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/patrykferenc/eecoin/internal/peer/command"
	"github.com/stretchr/testify/assert"
)

func TestAcceptPingController(t *testing.T) {
	// given
	acceptPingHandler := &noOpAcceptPingHandler{}
	pingController := NewPingController(acceptPingHandler)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// when
	pingController.GetPing(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Body.String())
}

type noOpAcceptPingHandler struct{}

func (h *noOpAcceptPingHandler) Handle(cmd command.AcceptPing) error {
	return nil
}
