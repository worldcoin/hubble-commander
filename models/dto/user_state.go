package dto

import "github.com/Worldcoin/hubble-commander/models"

type UserStateWithID struct {
	StateID uint32
	models.UserState
}

func MakeUserStateWithID(stateLeaf *models.StateLeaf) UserStateWithID {
	return UserStateWithID{
		StateID:   stateLeaf.StateID,
		UserState: stateLeaf.UserState,
	}
}

func NewUserStateWithID(stateLeaf *models.StateLeaf) *UserStateWithID {
	userState := MakeUserStateWithID(stateLeaf)
	return &userState
}
