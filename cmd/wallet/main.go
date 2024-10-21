package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/patrykferenc/eecoin/internal/hello"
	"github.com/patrykferenc/eecoin/internal/node/domain/node"
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
)

func main() {
	fmt.Printf("Wallet: %s\n", hello.Hello{Message: "wallet"}.String())

	sendTransactionToSetPeer()
}

func sendTransactionToSetPeer() {
	transaction := node.Transaction{
		ID:        blockchain.TransactionID(uuid.New().String()),
		Content:   "czemu taki poważny",
		Timestamp: time.Now(),
		From:      wallet.ID("mr.jajeczko"),
		To:        wallet.ID("miała-matka-syna"),
	}

	b, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transaction: %s\n", string(b))

	peer := "http://localhost:22137" + "/client/message"
	req, err := http.NewRequest(http.MethodPost, peer, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n", resp.Status)
}
