// nolint:lll
package config

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
)

func TestReadGenesisFile(t *testing.T) {
	genesisAccounts, err := readGenesisFile(getGenesisPath())
	require.NoError(t, err)
	require.Greater(t, len(genesisAccounts), 0)
	require.Equal(t, genesisAccounts[0].Balance.CmpN(0), 1)
	require.Equal(t, genesisAccounts[len(genesisAccounts)-1].Balance.CmpN(0), 0)
}

func TestDecodeRawGenesisAccounts(t *testing.T) {
	rawGenesisAccounts := []models.RawGenesisAccount{
		{
			PrivateKey: "2f7a559b2d2d4ec1e3babc0122e7ef0c6a45cdb4ccd167f456caca521123fe9e",
			Balance:    models.MakeUint256(1024),
		},
		{
			PrivateKey: "0131b1f02f2504a60a30261fa3665ca12ed542ab01f73debb19bafa996790272",
			Balance:    models.MakeUint256(0),
		},
	}

	expected := []models.GenesisAccount{
		{
			PrivateKey: [32]byte{47, 122, 85, 155, 45, 45, 78, 193, 227, 186, 188, 1, 34, 231, 239, 12, 106, 69, 205, 180, 204, 209, 103, 244, 86, 202, 202, 82, 17, 35, 254, 158},
			Balance:    models.MakeUint256(1024),
		},
		{
			PrivateKey: [32]byte{1, 49, 177, 240, 47, 37, 4, 166, 10, 48, 38, 31, 163, 102, 92, 161, 46, 213, 66, 171, 1, 247, 61, 235, 177, 155, 175, 169, 150, 121, 2, 114},
			Balance:    models.MakeUint256(0),
		},
	}

	genesisAccounts, err := decodeRawGenesisAccounts(rawGenesisAccounts)
	require.NoError(t, err)
	require.Equal(t, expected, genesisAccounts)
}
