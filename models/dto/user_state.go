package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type UserState struct {
	PubKeyID uint32
	TokenID  models.Uint256
	Balance  models.Uint256
	Nonce    models.Uint256
}

type UserStateWithID struct {
	StateID uint32
	UserState
}

func MakeUserStateWithID(stateID uint32, userState *models.UserState) UserStateWithID {
	return UserStateWithID{
		StateID: stateID,
		UserState: MakeUserState(userState),
	}
}

func MakeUserState(userState *models.UserState) UserState {
	return UserState{
		PubKeyID: userState.PubKeyID,
		TokenID:  userState.TokenID,
		Balance:  userState.Balance,
		Nonce:    userState.Nonce,
	}
}
