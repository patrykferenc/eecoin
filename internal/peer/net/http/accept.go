package http

import (
	"net/http"

	"github.com/patrykferenc/eecoin/internal/peer/command"
)

type PingController struct {
	acceptPingHandler command.AcceptPingHandler
}

func NewPingController(acceptPingHandler command.AcceptPingHandler) *PingController {
	return &PingController{acceptPingHandler: acceptPingHandler}
}

func (p *PingController) GetPing(w http.ResponseWriter, r *http.Request) {
	hostIP := r.RemoteAddr
	if err := p.acceptPingHandler.Handle(command.AcceptPing{Host: hostIP}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
