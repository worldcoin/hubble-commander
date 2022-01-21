package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadGenesisFile(t *testing.T) {
	genesisAccounts, err := readGenesisFile(getGenesisPath())
	require.NoError(t, err)
	require.Greater(t, len(genesisAccounts), 0)
	require.Equal(t, genesisAccounts[0].State.Balance.CmpN(0), 1)
}
