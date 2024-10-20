package node

import "github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"

type SendMessageEvent struct {
	TransactionID blockchain.TransactionID
}

type MessageSentEvent struct {
	TransactionID blockchain.TransactionID
}
