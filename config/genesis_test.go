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
	require.Equal(t, genesisAccounts[0].State.Balance.CmpN(0), 1)
}

func TestDecodeRawGenesisAccounts(t *testing.T) {
	rawGenesisAccounts := []models.RawGenesisAccount{
		{
			PublicKey: "26d15404bc8fafa97ea0b383e2bf17b7c2d34a9caabfa06eca5f748bb78e9d4128dc860d7644e3e6d8de3fe2c7d4f7097bc59c26b298e5029cfd37e78ffbfea41f92f018eaa32e1ad8c78861bb94b728be49703ff8ed2c7f46883cac26a73bec1255d71d0f906d83d079f0c8451322ea9e050f3b9143f09f8f81c8c53d69ced3",
			State: models.GenesisState{
				Balance: 1024,
			},
		},
		{
			PublicKey: "2c63e87bfc418b501b438fadf16f4a178a267ade00775d37443f15ebb376201c1c13ba71fb499492895a9eb0449bebcaeda1069641c7a95273b1f2ff7baaea780882ba670aef94491025b8a3ab3ca5fd352a9a2cd5e82f67bdb7164e8968d4d902f7ad7d42c8f325a73c2851b036dfe3282de47da444a0702d509e1a19058954",
		},
		{
			PublicKey: "26d15404bc8fafa97ea0b383e2bf17b7c2d34a9caabfa06eca5f748bb78e9d4128dc860d7644e3e6d8de3fe2c7d4f7097bc59c26b298e5029cfd37e78ffbfea41f92f018eaa32e1ad8c78861bb94b728be49703ff8ed2c7f46883cac26a73bec1255d71d0f906d83d079f0c8451322ea9e050f3b9143f09f8f81c8c53d69ced3",
			State: models.GenesisState{
				StateID:  1,
				PubKeyID: 2,
				TokenID:  3,
				Balance:  1024,
				Nonce:    4,
			},
		},
	}

	expected := []models.GenesisAccount{
		{
			PublicKey: models.PublicKey{38, 209, 84, 4, 188, 143, 175, 169, 126, 160, 179, 131, 226, 191, 23, 183, 194, 211, 74, 156, 170, 191, 160, 110, 202, 95, 116, 139, 183, 142, 157, 65, 40, 220, 134, 13, 118, 68, 227, 230, 216, 222, 63, 226, 199, 212, 247, 9, 123, 197, 156, 38, 178, 152, 229, 2, 156, 253, 55, 231, 143, 251, 254, 164, 31, 146, 240, 24, 234, 163, 46, 26, 216, 199, 136, 97, 187, 148, 183, 40, 190, 73, 112, 63, 248, 237, 44, 127, 70, 136, 60, 172, 38, 167, 59, 236, 18, 85, 215, 29, 15, 144, 109, 131, 208, 121, 240, 200, 69, 19, 34, 234, 158, 5, 15, 59, 145, 67, 240, 159, 143, 129, 200, 197, 61, 105, 206, 211},
			State: &models.StateLeaf{
				StateID: 0,
				UserState: models.UserState{
					PubKeyID: 0,
					TokenID:  models.MakeUint256(0),
					Balance:  models.MakeUint256(1024),
					Nonce:    models.MakeUint256(0),
				},
			},
		},
		{
			PublicKey: models.PublicKey{44, 99, 232, 123, 252, 65, 139, 80, 27, 67, 143, 173, 241, 111, 74, 23, 138, 38, 122, 222, 0, 119, 93, 55, 68, 63, 21, 235, 179, 118, 32, 28, 28, 19, 186, 113, 251, 73, 148, 146, 137, 90, 158, 176, 68, 155, 235, 202, 237, 161, 6, 150, 65, 199, 169, 82, 115, 177, 242, 255, 123, 170, 234, 120, 8, 130, 186, 103, 10, 239, 148, 73, 16, 37, 184, 163, 171, 60, 165, 253, 53, 42, 154, 44, 213, 232, 47, 103, 189, 183, 22, 78, 137, 104, 212, 217, 2, 247, 173, 125, 66, 200, 243, 37, 167, 60, 40, 81, 176, 54, 223, 227, 40, 45, 228, 125, 164, 68, 160, 112, 45, 80, 158, 26, 25, 5, 137, 84},
			State:     &models.StateLeaf{},
		},
		{
			PublicKey: models.PublicKey{38, 209, 84, 4, 188, 143, 175, 169, 126, 160, 179, 131, 226, 191, 23, 183, 194, 211, 74, 156, 170, 191, 160, 110, 202, 95, 116, 139, 183, 142, 157, 65, 40, 220, 134, 13, 118, 68, 227, 230, 216, 222, 63, 226, 199, 212, 247, 9, 123, 197, 156, 38, 178, 152, 229, 2, 156, 253, 55, 231, 143, 251, 254, 164, 31, 146, 240, 24, 234, 163, 46, 26, 216, 199, 136, 97, 187, 148, 183, 40, 190, 73, 112, 63, 248, 237, 44, 127, 70, 136, 60, 172, 38, 167, 59, 236, 18, 85, 215, 29, 15, 144, 109, 131, 208, 121, 240, 200, 69, 19, 34, 234, 158, 5, 15, 59, 145, 67, 240, 159, 143, 129, 200, 197, 61, 105, 206, 211},
			State: &models.StateLeaf{
				StateID: 1,
				UserState: models.UserState{
					PubKeyID: 2,
					TokenID:  models.MakeUint256(3),
					Balance:  models.MakeUint256(1024),
					Nonce:    models.MakeUint256(4),
				},
			},
		},
	}

	genesisAccounts, err := decodeRawGenesisAccounts(rawGenesisAccounts)
	require.NoError(t, err)
	require.Equal(t, expected, genesisAccounts)
}

func TestDecodeRawGenesisAccounts_InvalidPublicKeyLength(t *testing.T) {
	rawGenesisAccounts := []models.RawGenesisAccount{
		{
			PublicKey: "26d1",
		},
	}

	_, err := decodeRawGenesisAccounts(rawGenesisAccounts)
	require.Error(t, err)
	require.Equal(t, models.ErrInvalidPublicKeyLength, err)
}

func TestTestDecodeRawGenesisAccounts_MissingPublicKey(t *testing.T) {
	matchingKeys := []models.RawGenesisAccount{
		{
			State: models.GenesisState{
				StateID: 1,
				Balance: 100,
			},
		},
	}
	_, err := decodeRawGenesisAccounts(matchingKeys)
	require.ErrorIs(t, err, errMissingGenesisPublicKey)
}
