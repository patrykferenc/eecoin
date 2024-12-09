package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnspentClient_Get(t *testing.T) {
	remote := "http://localhost:8080"
	client := NewUnspentOutputsRepository(remote)
	assert.NotNil(t, client)
}
