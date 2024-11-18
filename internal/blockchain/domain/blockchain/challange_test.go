package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChallange_shouldError(t *testing.T) {
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
			description: "Difficulty too big and not even",
			difficulty:  65,
		},
		{
			description: "Difficulty too big",
			difficulty:  66,
		},
		{
			description: "Difficulty not even",
			difficulty:  3,
		},
	}

	for _, tc := range tt {
		// when
		_, err := NewChallenge(tc.difficulty)

		// then
		assertThat.Equal(NotValidDifficulty, err)
	}
}

func TestNewChallange_shouldWork(t *testing.T) {
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
			description: "Difficulty 64",
			difficulty:  64,
		},
		{
			description: "Difficulty 4",
			difficulty:  4,
		},
	}

	for _, tc := range tt {
		// when
		challenge, err := NewChallenge(tc.difficulty)

		// then
		assertThat.Nil(err)
		assertThat.Equal(tc.difficulty, challenge.Difficulty)
	}
}

func TestRollNonce(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	challenge, err := NewChallenge(2)
	assertThat.Nil(err)

	// and given
	initialNonce := challenge.Nonce
	initialHash := challenge.HashValue

	// when
	challenge.RollNonce()

	// then
	assertThat.NotEqual(initialNonce, challenge.Nonce)
	assertThat.NotEqual(initialHash, challenge.HashValue)
}

func TestChallengeMatchesDifficulty_shouldMatch(t *testing.T) {
	t.Parallel()
	assertThat := assert.New(t)

	// given
	challenge := Challenge{Difficulty: 4, HashValue: []byte{0, 0, 1}}

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
			challenge:   Challenge{Difficulty: 4, HashValue: []byte{0, 1, 1}},
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
	challenge, _ := NewChallenge(2)
	initialNonce := challenge.Nonce
	initialHash := challenge.HashValue

	// when
	challenge.RollUntilMatchesDifficultyCapped(2)

	// then
	assertThat.NotEqual(initialNonce, challenge.Nonce)
	assertThat.NotEqual(initialHash, challenge.HashValue)
}

func mustNewChallenge(difficulty int) Challenge {
	challenge, err := NewChallenge(difficulty)
	if err != nil {
		panic(err)
	}
	return challenge
}
