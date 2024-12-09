package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnspentClient_Get(t *testing.T) {
	c := http.Client{}
	remote := "http://localhost:8080"
	client := NewUnspentOutputsRepository(c, remote)
	assert.NotNil(t, client)
}
