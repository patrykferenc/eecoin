package persistence

import (
	"testing"

	"github.com/patrykferenc/eecoin/internal/blockchain/domain/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersistBlockChain(t *testing.T) {
	assertThat := assert.New(t)

	// given
	dir := t.TempDir() + "/chain"
	genesis := blockchain.GenerateGenesisBlock()
	chain, err := blockchain.ImportBlockchain([]blockchain.Block{genesis})
	require.NoError(t, err)

	// when - then
	saveErr := Persist(*chain, dir)
	assertThat.Nil(saveErr)

	// when - then
	loaded, err := Load(dir)

	assertThat.Nil(err)
	assertThat.NotNil(loaded)
	assertThat.Equal(loaded.GetFirst(), chain.GetFirst())
	assertThat.Equal(loaded.GetLast(), chain.GetLast())
}
