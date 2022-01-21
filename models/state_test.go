package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestUserState_SetBytes(t *testing.T) {
	userState := UserState{
		PubKeyID: 3,
		TokenID:  MakeUint256(4),
		Balance:  MakeUint256(5),
		Nonce:    MakeUint256(6),
	}

	var decodedUserState UserState
	_ = decodedUserState.SetBytes(userState.Bytes())
	require.Equal(t, userState, decodedUserState)
}

func TestUserState_SetBytes_InvalidLength(t *testing.T) {
	bytes := make([]byte, 18)
	userState := UserState{}
	err := userState.SetBytes(bytes)
	require.ErrorIs(t, err, ErrInvalidLength)
}

func TestUserState_Copy(t *testing.T) {
	userState := UserState{
		PubKeyID: 1,
		TokenID:  MakeUint256(2),
		Balance:  MakeUint256(3),
		Nonce:    MakeUint256(4),
	}
	userStateCopy := userState.Copy()
	userStateCopy.Balance = *userStateCopy.Balance.AddN(100)
	require.Equal(t, MakeUint256(3), userState.Balance)
}

func TestStateUpdate_SetBytes(t *testing.T) {
	stateUpdate := StateUpdate{
		ID:          1,
		CurrentRoot: utils.RandomHash(),
		PrevRoot:    utils.RandomHash(),
		PrevStateLeaf: StateLeaf{
			StateID:  2,
			DataHash: utils.RandomHash(),
			UserState: UserState{
				PubKeyID: 3,
				TokenID:  MakeUint256(4),
				Balance:  MakeUint256(5),
				Nonce:    MakeUint256(6),
			},
		},
	}

	var decodedUpdate StateUpdate
	_ = decodedUpdate.SetBytes(stateUpdate.Bytes())
	require.Equal(t, stateUpdate, decodedUpdate)
}

func TestStateUpdate_SetBytes_InvalidLength(t *testing.T) {
	bytes := make([]byte, 18)
	stateUpdate := StateUpdate{}
	err := stateUpdate.SetBytes(bytes)
	require.ErrorIs(t, err, ErrInvalidLength)
}
