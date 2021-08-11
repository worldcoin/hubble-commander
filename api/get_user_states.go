package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]dto.UserStateWithID, error) {
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
