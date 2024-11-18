package blockchain

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand/v2"
)

var NotValidDifficulty = errors.New("difficulty not valid, it must divisible by 2 and between 2 and 64")

type Challenge struct {
	Difficulty int
	Nonce      uint32
	HashValue  []byte
}

func (c *Challenge) RollNonce() {
	newNonce := rand.Uint32()

	nonceByteBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceByteBuffer, newNonce)
	nonceHash := sha256.Sum256(nonceByteBuffer)

	c.HashValue = nonceHash[:]
	c.Nonce = newNonce
}

func (c *Challenge) MatchesDifficulty() bool {
	if c.HashValue == nil {
		return false
	}
	for i := 0; i < c.Difficulty/2; i++ {
		if c.HashValue[i] != 0 {
			return false
		}
	}
	return true
}

func (c *Challenge) RollUntilMatchesDifficulty() {
	for i := 0; c.MatchesDifficulty(); i++ {
		c.RollNonce()
	}
}

func (c *Challenge) RollUntilMatchesDifficultyCapped(maxIterations int) {
	for i := 0; i < maxIterations || c.MatchesDifficulty(); i++ {
		c.RollNonce()
	}
}

func NewChallenge(difficulty int) (Challenge, error) {
	if difficulty >= 2 && difficulty <= 64 && difficulty%2 == 0 {
		return Challenge{
			Difficulty: difficulty,
		}, nil
	}
	return Challenge{}, NotValidDifficulty
}
