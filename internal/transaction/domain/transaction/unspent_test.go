package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnspentWhenValid(t *testing.T) {
	assert := assert.New(t)
	// given
	someUnspents := []UnspentOutput{NewUnspentOutput("someID-1", 0, 100, "someAddress-1"), NewUnspentOutput("someID-2", 1, 200, "someAddress-2"), NewUnspentOutput("someID-3", 2, 300, "someAddress-3")}

	// when
	leftover, included, err := calculateUnspentForAmount(someUnspents, 300)

	// then
	assert.NoError(err)

	// and then
	assert.Equal(0, leftover)
	assert.Len(included, 2) // since there is no ordering, we take the first ones
}

func TestUnspentWhenNotEnough(t *testing.T) {
	assert := assert.New(t)
	// given
	someUnspents := []UnspentOutput{NewUnspentOutput("someID-1", 0, 100, "someAddress-1"), NewUnspentOutput("someID-2", 1, 200, "someAddress-2")}

	// when
	_, _, err := calculateUnspentForAmount(someUnspents, 600)

	// then
	assert.Error(err)
}
