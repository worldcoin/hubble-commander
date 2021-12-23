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

func MakeUserStateWithID(stateLeaf *models.StateLeaf) UserStateWithID {
	return UserStateWithID{
		StateID: stateLeaf.StateID,
		UserState: UserState{
			PubKeyID: stateLeaf.UserState.PubKeyID,
			TokenID:  stateLeaf.UserState.TokenID,
			Balance:  stateLeaf.UserState.Balance,
			Nonce:    stateLeaf.UserState.Nonce,
		},
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

func NewUserStateWithID(stateLeaf *models.StateLeaf) *UserStateWithID {
	userState := MakeUserStateWithID(stateLeaf)
	return &userState
}
