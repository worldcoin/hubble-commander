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
	require.Equal(t, genesisAccounts[len(genesisAccounts)-1].State.Balance.CmpN(0), 0)
}

func TestDecodeRawGenesisAccounts(t *testing.T) {
	rawGenesisAccounts := []models.RawGenesisAccount{
		{
			PublicKey:  "26d15404bc8fafa97ea0b383e2bf17b7c2d34a9caabfa06eca5f748bb78e9d4128dc860d7644e3e6d8de3fe2c7d4f7097bc59c26b298e5029cfd37e78ffbfea41f92f018eaa32e1ad8c78861bb94b728be49703ff8ed2c7f46883cac26a73bec1255d71d0f906d83d079f0c8451322ea9e050f3b9143f09f8f81c8c53d69ced3",
			PrivateKey: "01e6cf81f2726600430b08581b39a431fa847250519a8c608ca57379a21883f2",
			State: models.GenesisState{
				Balance: 1024,
			},
		},
		{
			PublicKey:  "2c63e87bfc418b501b438fadf16f4a178a267ade00775d37443f15ebb376201c1c13ba71fb499492895a9eb0449bebcaeda1069641c7a95273b1f2ff7baaea780882ba670aef94491025b8a3ab3ca5fd352a9a2cd5e82f67bdb7164e8968d4d902f7ad7d42c8f325a73c2851b036dfe3282de47da444a0702d509e1a19058954",
			PrivateKey: "02393c7d3d803f1b257f727c58a7b85f5e9283c64439ec9a4efb40acf9b6a5af",
			Balance:    0,
		},
		{
			PublicKey: "0cb8f5ab36949be73fee1e99d209fd6bc7d7b3d3150ab9834fca795a474163a2155311fe991035c80f48a8416f841f8b7d885424a85e994aebc0c5ce82564aa91a32e8c98d7b9808212cf4d661ca03d6e49e45fdce1f1bbb14291a8d8b2743eb0099cc354c7b39a8625d63ec344c6210b941d27d0cb4805d498dd97c19b2e8b2",
			Balance:   0,
		},
		{
			PrivateKey: "2bd0e28569bcbe3306186d81085fa595f3a716243d02c3ae301538bf91f93b50",
			Balance:    0,
		},
		{
			PublicKey:  "26d15404bc8fafa97ea0b383e2bf17b7c2d34a9caabfa06eca5f748bb78e9d4128dc860d7644e3e6d8de3fe2c7d4f7097bc59c26b298e5029cfd37e78ffbfea41f92f018eaa32e1ad8c78861bb94b728be49703ff8ed2c7f46883cac26a73bec1255d71d0f906d83d079f0c8451322ea9e050f3b9143f09f8f81c8c53d69ced3",
			PrivateKey: "01e6cf81f2726600430b08581b39a431fa847250519a8c608ca57379a21883f2",
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
			PublicKey:  &models.PublicKey{38, 209, 84, 4, 188, 143, 175, 169, 126, 160, 179, 131, 226, 191, 23, 183, 194, 211, 74, 156, 170, 191, 160, 110, 202, 95, 116, 139, 183, 142, 157, 65, 40, 220, 134, 13, 118, 68, 227, 230, 216, 222, 63, 226, 199, 212, 247, 9, 123, 197, 156, 38, 178, 152, 229, 2, 156, 253, 55, 231, 143, 251, 254, 164, 31, 146, 240, 24, 234, 163, 46, 26, 216, 199, 136, 97, 187, 148, 183, 40, 190, 73, 112, 63, 248, 237, 44, 127, 70, 136, 60, 172, 38, 167, 59, 236, 18, 85, 215, 29, 15, 144, 109, 131, 208, 121, 240, 200, 69, 19, 34, 234, 158, 5, 15, 59, 145, 67, 240, 159, 143, 129, 200, 197, 61, 105, 206, 211},
			PrivateKey: &[32]byte{1, 230, 207, 129, 242, 114, 102, 0, 67, 11, 8, 88, 27, 57, 164, 49, 250, 132, 114, 80, 81, 154, 140, 96, 140, 165, 115, 121, 162, 24, 131, 242},
			Balance:    models.MakeUint256(0),
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
			PublicKey:  &models.PublicKey{44, 99, 232, 123, 252, 65, 139, 80, 27, 67, 143, 173, 241, 111, 74, 23, 138, 38, 122, 222, 0, 119, 93, 55, 68, 63, 21, 235, 179, 118, 32, 28, 28, 19, 186, 113, 251, 73, 148, 146, 137, 90, 158, 176, 68, 155, 235, 202, 237, 161, 6, 150, 65, 199, 169, 82, 115, 177, 242, 255, 123, 170, 234, 120, 8, 130, 186, 103, 10, 239, 148, 73, 16, 37, 184, 163, 171, 60, 165, 253, 53, 42, 154, 44, 213, 232, 47, 103, 189, 183, 22, 78, 137, 104, 212, 217, 2, 247, 173, 125, 66, 200, 243, 37, 167, 60, 40, 81, 176, 54, 223, 227, 40, 45, 228, 125, 164, 68, 160, 112, 45, 80, 158, 26, 25, 5, 137, 84},
			PrivateKey: &[32]byte{2, 57, 60, 125, 61, 128, 63, 27, 37, 127, 114, 124, 88, 167, 184, 95, 94, 146, 131, 198, 68, 57, 236, 154, 78, 251, 64, 172, 249, 182, 165, 175},
			Balance:    models.MakeUint256(0),
			State:      &models.StateLeaf{},
		},
		{
			PublicKey:  &models.PublicKey{12, 184, 245, 171, 54, 148, 155, 231, 63, 238, 30, 153, 210, 9, 253, 107, 199, 215, 179, 211, 21, 10, 185, 131, 79, 202, 121, 90, 71, 65, 99, 162, 21, 83, 17, 254, 153, 16, 53, 200, 15, 72, 168, 65, 111, 132, 31, 139, 125, 136, 84, 36, 168, 94, 153, 74, 235, 192, 197, 206, 130, 86, 74, 169, 26, 50, 232, 201, 141, 123, 152, 8, 33, 44, 244, 214, 97, 202, 3, 214, 228, 158, 69, 253, 206, 31, 27, 187, 20, 41, 26, 141, 139, 39, 67, 235, 0, 153, 204, 53, 76, 123, 57, 168, 98, 93, 99, 236, 52, 76, 98, 16, 185, 65, 210, 125, 12, 180, 128, 93, 73, 141, 217, 124, 25, 178, 232, 178},
			PrivateKey: nil,
			Balance:    models.MakeUint256(0),
			State:      &models.StateLeaf{},
		},
		{
			PublicKey:  nil,
			PrivateKey: &[32]byte{43, 208, 226, 133, 105, 188, 190, 51, 6, 24, 109, 129, 8, 95, 165, 149, 243, 167, 22, 36, 61, 2, 195, 174, 48, 21, 56, 191, 145, 249, 59, 80},
			Balance:    models.MakeUint256(0),
			State:      &models.StateLeaf{},
		},
		{
			PublicKey:  &models.PublicKey{38, 209, 84, 4, 188, 143, 175, 169, 126, 160, 179, 131, 226, 191, 23, 183, 194, 211, 74, 156, 170, 191, 160, 110, 202, 95, 116, 139, 183, 142, 157, 65, 40, 220, 134, 13, 118, 68, 227, 230, 216, 222, 63, 226, 199, 212, 247, 9, 123, 197, 156, 38, 178, 152, 229, 2, 156, 253, 55, 231, 143, 251, 254, 164, 31, 146, 240, 24, 234, 163, 46, 26, 216, 199, 136, 97, 187, 148, 183, 40, 190, 73, 112, 63, 248, 237, 44, 127, 70, 136, 60, 172, 38, 167, 59, 236, 18, 85, 215, 29, 15, 144, 109, 131, 208, 121, 240, 200, 69, 19, 34, 234, 158, 5, 15, 59, 145, 67, 240, 159, 143, 129, 200, 197, 61, 105, 206, 211},
			PrivateKey: &[32]byte{1, 230, 207, 129, 242, 114, 102, 0, 67, 11, 8, 88, 27, 57, 164, 49, 250, 132, 114, 80, 81, 154, 140, 96, 140, 165, 115, 121, 162, 24, 131, 242},
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
			PublicKey:  "26d1",
			PrivateKey: "",
			Balance:    1024,
		},
	}

	_, err := decodeRawGenesisAccounts(rawGenesisAccounts)
	require.Error(t, err)
	require.Equal(t, models.ErrInvalidPublicKeyLength, err)
}

func TestTestDecodeRawGenesisAccounts_ValidatesThatKeysMatch(t *testing.T) {
	matchingKeys := []models.RawGenesisAccount{
		{
			PublicKey:  "0df68cb87856229b0bc3f158fff8b82b04deb1a4c23dadbf3ed2da4ec6f6efcb1c165c6b47d8c89ab2ddb0831c182237b27a4b3d9701775ad6c180303f87ef260566cb2f0bcc7b89c2260de2fee8ec29d7b5e575a1e36eb4bcead52a74a511b7188d7df7c9d08f94b9daa9d89105fbdf22bf14e30b84f8adefb3695ebff00e88",
			PrivateKey: "2f7a559b2d2d4ec1e3babc0122e7ef0c6a45cdb4ccd167f456caca521123fe9e",
			Balance:    1024,
		},
	}
	_, err := decodeRawGenesisAccounts(matchingKeys)
	require.NoError(t, err)

	nonMatchingKeys := []models.RawGenesisAccount{
		{
			PublicKey:  "0df68cb87856229b0bc3f158fff8b82b04deb1a4c23dadbf3ed2da4ec6f6efcb1c165c6b47d8c89ab2ddb0831c182237b27a4b3d9701775ad6c180303f87ef260566cb2f0bcc7b89c2260de2fee8ec29d7b5e575a1e36eb4bcead52a74a511b7188d7df7c9d08f94b9daa9d89105fbdf22bf14e30b84f8adefb3695ebff00e88",
			PrivateKey: "0131b1f02f2504a60a30261fa3665ca12ed542ab01f73debb19bafa996790272",
			Balance:    1024,
		},
	}
	_, err = decodeRawGenesisAccounts(nonMatchingKeys)
	require.ErrorIs(t, err, NewErrNonMatchingKeys("0x"+nonMatchingKeys[0].PublicKey))
}
