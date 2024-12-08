package blockchain

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChallenge_shouldError(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	tt := []struct {
		description string
		difficulty  int
	}{
		{
			description: "Difficulty too small",
			difficulty:  1,
		},
		{
			description: "Difficulty too big",
			difficulty:  257,
		},
	}

	for _, tc := range tt {
		// when
		_, err := NewChallenge(tc.difficulty, 60)

		// then
		assertThat.Equal(NotValidDifficulty, err)
	}
}

func TestNewChallenge_shouldWork(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	tt := []struct {
		description string
		difficulty  int
	}{
		{
			description: "Difficulty 2",
			difficulty:  2,
		},
		{
			description: "Difficulty 32",
			difficulty:  32,
		},
		{
			description: "Difficulty 4",
			difficulty:  4,
		},
	}

	for _, tc := range tt {
		// when
		challenge, err := NewChallenge(tc.difficulty, 60)

		// then
		assertThat.Nil(err)
		assertThat.Equal(tc.difficulty, challenge.Difficulty)
	}
}

func TestRollNonce(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	challenge, err := NewChallenge(2, 60)
	assertThat.Nil(err)

	// and given
	initialNonce := challenge.Nonce
	initialHash := challenge.HashValue

	// and given
	prePreparedBlock := GenerateGenesisBlock()

	// when
	err = challenge.RollNonce(prePreparedBlock, 60)

	// then
	assertThat.Nil(err)
	assertThat.NotEqual(initialNonce, challenge.Nonce)
	assertThat.NotEqual(initialHash, challenge.HashValue)
}

func TestChallengeMatchesDifficulty_shouldMatch(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	challangeHashStr := base64.StdEncoding.EncodeToString([]byte{0})
	challenge := Challenge{Difficulty: 4, HashValue: challangeHashStr}

	// then
	assertThat.True(challenge.MatchesDifficulty())
}

func TestChallengeMatchesDifficulty_shouldNotMatch(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	tt := []struct {
		description string
		challenge   Challenge
	}{
		{
			description: "Invalid hash",
			challenge:   Challenge{Difficulty: 4, HashValue: "01"},
		},
		{
			description: "Nil hash",
			challenge:   mustNewChallenge(4),
		},
	}

	for _, tc := range tt {
		// then
		assertThat.False(tc.challenge.MatchesDifficulty())
	}
}

func TestRollUntil(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	challenge, _ := NewChallenge(2, 60)
	initialNonce := challenge.Nonce
	initialHash := challenge.HashValue

	// and given
	prePreparedBlock := GenerateGenesisBlock()

	// when
	err := challenge.RollUntilMatchesDifficultyCapped(2, prePreparedBlock, 60)

	// then
	assertThat.Nil(err)
	assertThat.NotEqual(initialNonce, challenge.Nonce)
	assertThat.NotEqual(initialHash, challenge.HashValue)
}

func mustNewChallenge(difficulty int) Challenge {
	challenge, err := NewChallenge(difficulty, 60)
	if err != nil {
		panic(err)
	}
	return challenge
}
