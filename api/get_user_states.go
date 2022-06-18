package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getUserStatesAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(99003, "user states not found"),
}

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]dto.UserStateWithID, error) {
	batch, err := a.unsafeGetUserStates(publicKey)
	if err != nil {
		return nil, sanitizeError(err, getUserStatesAPIErrors)
	}

	return batch, nil
}

func (a *API) unsafeGetUserStates(publicKey *models.PublicKey) ([]dto.UserStateWithID, error) {
	// span
	leaves, err := a.storage.GetStateLeavesByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	userStates := make([]dto.UserStateWithID, 0, len(leaves))
	for i := range leaves {
		userStates = append(userStates, dto.MakeUserStateWithID(&leaves[i]))
	}

	return userStates, nil
}
