package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPopulatedGenesisAccount_Bytes_ReturnsACopy(t *testing.T) {
	account := PopulatedGenesisAccount{
		PublicKey: PublicKey{1, 2, 0, 5, 4},
	}
	bytes := account.Bytes()
	bytes[0] = 9
	require.Equal(t, PublicKey{1, 2, 0, 5, 4}, account.PublicKey)
}

func TestPopulatedGenesisAccount_SetBytes(t *testing.T) {
	account := PopulatedGenesisAccount{
		PublicKey: PublicKey{1, 2, 0, 5, 4},
		StateID:   44,
		State: UserState{
			PubKeyID: 7,
			TokenID:  MakeUint256(0),
			Balance:  MakeUint256(4314),
			Nonce:    MakeUint256(0),
		},
	}
	bytes := account.Bytes()
	newAccount := PopulatedGenesisAccount{}
	err := newAccount.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, account, newAccount)
}
