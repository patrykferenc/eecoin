package main

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/hello"
)

func main() {
	fmt.Printf("Node: %s\n", hello.Hello{Message: "node"}.String())
}
