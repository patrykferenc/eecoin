package main

import (
	"fmt"

	"github.com/patrykferenc/eecoin/internal/hello"
)

func main() {
	fmt.Printf("Node: %s\n", hello.Hello{Message: "node"}.String())

	// TODO: add handlers of events
	// create components
	// open the file for reading da peers from da hood
	// start the goddamn server with all the routs and stuff : >
}
