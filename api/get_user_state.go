package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getUserStateAPIErrors = map[error]ErrorAPI{
	&storage.NotFoundError{}: {
		Code:    10000,
		Message: "user state not found",
	},
}

func (a *API) GetUserState(id uint32) (*dto.UserStateWithID, error) {
	userState, err := a.unsafeGetUserState(id)
	if err != nil {
		return nil, sanitizeError(err, getUserStateAPIErrors)
	}

	return userState, nil
}

func (a *API) unsafeGetUserState(id uint32) (*dto.UserStateWithID, error) {
	leaf, err := a.storage.StateTree.Leaf(id)
	if err != nil {
		return nil, err
	}
	return dto.NewUserStateWithID(leaf), nil
}
