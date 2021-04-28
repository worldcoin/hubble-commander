package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGenesisAccounts(t *testing.T) {
	genesisAccounts, err := getGenesisAccounts("genesis.yaml")
	require.NoError(t, err)
	require.Greater(t, len(genesisAccounts), 0)
	require.Equal(t, genesisAccounts[0].Balance.CmpN(0), 1)
	require.Equal(t, genesisAccounts[len(genesisAccounts) - 1].Balance.CmpN(0), 0)
}
