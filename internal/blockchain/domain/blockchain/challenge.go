package blockchain

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"github.com/gymshark/go-hasher"
	"math/bits"
	"math/rand/v2"
)

var (
	NotValidDifficulty = errors.New("difficulty not valid, it must be between 2 and 256")
)

type Challenge struct {
	Difficulty    int
	Nonce         uint32
	HashValue     string
	TimeCapMillis int64
}

func (c *Challenge) RollNonce(previousBlock Block, currentTimestampMillis int64) error {
	newNonce := rand.Uint32()
	targetHash, err := calculateTargetHash(previousBlock, currentTimestampMillis, newNonce)
	if err != nil {
		return err
	}
	c.HashValue = targetHash
	c.Nonce = newNonce
	return nil
}

func (c *Challenge) MatchesDifficulty() bool {
	byteVal, err := base64.StdEncoding.DecodeString(c.HashValue)
	if err != nil {
		return false
	}
	if len(byteVal) < c.Difficulty/8 || len(byteVal) == 0 {
		return false
	}
	for i := 0; i <= c.Difficulty/8; i++ {
		if bits.LeadingZeros8(byteVal[i]) < c.Difficulty-(8*i) {
			return false
		}
	}
	return true
}

func (c *Challenge) RollUntilMatchesDifficulty(previousBlock Block, currentTimestampMillis int64) error {
	for i := 0; !c.MatchesDifficulty(); i++ {
		err := c.RollNonce(previousBlock, currentTimestampMillis)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Challenge) RollUntilMatchesDifficultyCapped(maxIterations int, previousBlock Block, currentTimestampMillis int64) error {
	for i := 0; i < maxIterations || !c.MatchesDifficulty(); i++ {
		err := c.RollNonce(previousBlock, currentTimestampMillis)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewChallenge(difficulty int, timeCapMillis int64) (Challenge, error) {
	if difficulty >= 2 && difficulty <= 256 {
		return Challenge{
			Difficulty:    difficulty,
			TimeCapMillis: timeCapMillis,
		}, nil
	}
	return Challenge{}, NotValidDifficulty
}

func Verify(previous Block, latestBlockTimestamp int64, latestSolvedChallengeNonce uint32, latestSolvedChallengeHashValue string) bool {
	validHash, err := calculateTargetHash(previous, latestBlockTimestamp, latestSolvedChallengeNonce)
	if err != nil || validHash != latestSolvedChallengeHashValue {
		return false
	}
	return true
}

func calculateTargetHash(previousBlock Block, currentTimestampMillis int64, nonce uint32) (string, error) {
	var allBytes []byte
	nextIndex := previousBlock.Index + 1
	previousHash := previousBlock.ContentHash
	previousBlockData := previousBlock.Transactions

	nonceByteBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceByteBuffer, nonce)

	nextIndexByteBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint32(nextIndexByteBuffer, uint32(nextIndex))

	currentTimestampMillisBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(nextIndexByteBuffer, uint64(currentTimestampMillis))

	previousHashByteBuffer := []byte(previousHash)

	var previousBlockDataBuffer bytes.Buffer
	enc := gob.NewEncoder(&previousBlockDataBuffer)
	if err := enc.Encode(previousBlockData); err != nil {
		return "", err
	}

	allBytes = append(allBytes, nonceByteBuffer...)
	allBytes = append(allBytes, nextIndexByteBuffer...)
	allBytes = append(allBytes, previousHashByteBuffer...)
	allBytes = append(allBytes, currentTimestampMillisBuffer...)
	allBytes = append(allBytes, previousBlockDataBuffer.Bytes()...)

	return hasher.Sha256(allBytes).Base64(), nil
}
