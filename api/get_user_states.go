package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]dto.ReturnUserState, error) {
	states, err := a.storage.GetUserStatesByPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	userStates := make([]dto.ReturnUserState, 0, len(states))
	for i := range states {
		userState := dto.ReturnUserState{
			UserState: states[i].UserState,
			StateID:   states[i].MerklePath.Path,
		}
		userStates = append(userStates, userState)
	}

	return userStates, nil
}
