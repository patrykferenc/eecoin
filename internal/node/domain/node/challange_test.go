package node

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRollNonce(t *testing.T) {
	assertThat := assert.New(t)

	challenge, errorThatShouldBeNil := NewChallenge(2)
	_, errorThatShouldNotBeNil := NewChallenge(1)

	initialNonce := challenge.Nonce
	initialHash := challenge.HashValue

	challenge.RollNonce()

	assertThat.NotEqual(initialNonce, challenge.Nonce)
	assertThat.NotEqual(initialHash, challenge.HashValue)

	assertThat.Nil(errorThatShouldBeNil)
	assertThat.Equal(errorThatShouldNotBeNil, NotValidDifficulty)
}

func TestChallengeMatchesDifficulty(t *testing.T) {
	assertThat := assert.New(t)

	challengePreparedWithValidHash := Challenge{Difficulty: 4, HashValue: []byte{0, 0, 1}}
	challengePreparedWithInvalidHash := Challenge{Difficulty: 4, HashValue: []byte{0, 1, 1}}
	challengeNilHash, _ := NewChallenge(2)

	assertThat.True(challengePreparedWithValidHash.MatchesDifficulty())
	assertThat.False(challengePreparedWithInvalidHash.MatchesDifficulty())
	assertThat.False(challengeNilHash.MatchesDifficulty())

}

func TestRollUntil(t *testing.T) {
	assertThat := assert.New(t)

	challenge, _ := NewChallenge(2)
	initialNonce := challenge.Nonce
	initialHash := challenge.HashValue

	challenge.RollUntilMatchesDifficultyCapped(2)

	assertThat.NotEqual(initialNonce, challenge.Nonce)
	assertThat.NotEqual(initialHash, challenge.HashValue)

}
