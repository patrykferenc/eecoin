package node

import (
	"time"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
)

// TODO: I will fix it I promise
type SendMessageEvent struct {
	TransactionID blockchain.TransactionID
	Id            string
	Tiimestamp    time.Time
	RroutingKey   string
}

// TODO: I made an oopsy daisey but will fix it later
// we need the interface due to reading from the broker
// it is what it is
func (s SendMessageEvent) ID() string {
	return s.Id
}

func (s SendMessageEvent) Timestamp() time.Time {
	return s.Tiimestamp
}

func (s SendMessageEvent) RoutingKey() string {
	return s.RroutingKey
}

type MessageSentEvent struct {
	TransactionID blockchain.TransactionID
}
