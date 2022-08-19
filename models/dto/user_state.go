package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type UserState struct {
	PubKeyID int64
	TokenID  models.Uint256
	Balance  models.Uint256
	Nonce    models.Uint256
}

type UserStateWithID struct {
	StateID int64
	UserState
}

func MakeUserStateWithID(stateID uint32, userState *models.UserState) UserStateWithID {
	return UserStateWithID{
		StateID:   int64(stateID),
		UserState: MakeUserState(userState),
	}
}

func MakePendingUserStateWithID(userState *models.UserState) UserStateWithID {
	return UserStateWithID{
		StateID:   -1,
		UserState: MakePendingUserState(userState),
	}
}

func MakeUserState(userState *models.UserState) UserState {
	return UserState{
		PubKeyID: int64(userState.PubKeyID),
		TokenID:  userState.TokenID,
		Balance:  userState.Balance,
		Nonce:    userState.Nonce,
	}
}

func MakePendingUserState(userState *models.UserState) UserState {
	maxUint32 := ^uint32(0)

	var pubkeyToReturn int64
	if userState.PubKeyID == maxUint32 {
		pubkeyToReturn = -1
	} else {
		pubkeyToReturn = int64(userState.PubKeyID)
	}

	return UserState{
		PubKeyID: pubkeyToReturn,
		TokenID:  userState.TokenID,
		Balance:  userState.Balance,
		Nonce:    userState.Nonce,
	}
}
