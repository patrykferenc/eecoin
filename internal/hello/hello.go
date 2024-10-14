package hello

import "fmt"

type Hello struct {
	Message string
}

func (h Hello) String() string {
	return fmt.Sprintf("Hello, %s", h.Message)
}
