package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]dto.UserState, error) {
	states, err := a.storage.GetUserStatesByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	userStates := make([]dto.UserState, 0, len(states))
	for i := range states {
		userState := dto.UserState{
			UserState: states[i].UserState,
			StateID:   states[i].MerklePath.Path,
		}
		userStates = append(userStates, userState)
	}

	return userStates, nil
}
