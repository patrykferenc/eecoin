package main

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/hello"
)

func main() {
	fmt.Printf("Wallet: %s\n", hello.Hello{Message: "wallet"}.String())
}
