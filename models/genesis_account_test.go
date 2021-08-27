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
		PubKeyID:  7,
		StateID:   44,
		Balance:   MakeUint256(4314),
	}
	bytes := account.Bytes()
	newAccount := PopulatedGenesisAccount{}
	err := newAccount.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, account, newAccount)
}
