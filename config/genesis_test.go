package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadGenesisFile(t *testing.T) {
	genesisPath := path.Join("..", "genesis.yaml")
	genesisAccounts, err := readGenesisFile(genesisPath)
	require.NoError(t, err)
	require.Greater(t, len(genesisAccounts), 0)
	require.Equal(t, genesisAccounts[0].State.Balance.CmpN(0), 1)
}
