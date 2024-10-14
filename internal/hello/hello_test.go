package hello

import "testing"

func TestHello(t *testing.T) {
	h := Hello{Message: "world"}

	if h.String() != "Hello, world" {
		t.Fatal("unexpected string")
	}
}
