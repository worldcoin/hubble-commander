package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getUserStateAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(99002, "user state not found"),
}

func (a *API) GetUserState(id uint32) (*dto.UserStateWithID, error) {
	userState, err := a.unsafeGetUserState(id)
	if err != nil {
		return nil, sanitizeError(err, getUserStateAPIErrors)
	}

	return userState, nil
}

func (a *API) unsafeGetUserState(id uint32) (*dto.UserStateWithID, error) {
	// span
	leaf, err := a.storage.StateTree.Leaf(id)
	if err != nil {
		return nil, err
	}
	return dto.NewUserStateWithID(leaf), nil
}
