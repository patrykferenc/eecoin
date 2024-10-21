package http

import (
	"net/http"
	"strings"

	"github.com/patrykferenc/eecoin/internal/peer/command"
)

func getPing(acceptPingHandler command.AcceptPingHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hostIP := strings.Split(r.RemoteAddr, ":")[0]                                                             // TODO: handle ips better - maybe use X-Forwarded-For or get it from the request
		if err := acceptPingHandler.Handle(command.AcceptPing{Host: "http://" + hostIP + ":22137"}); err != nil { // TODO: make port configurable (and clean up the rest)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
