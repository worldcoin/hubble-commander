package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]models.ReturnUserState, error) {
	states, err := a.storage.GetUserStates(publicKey)
	if err != nil {
		return nil, err
	}

	userStates := make([]models.ReturnUserState, 0, len(states))
	for i := range states {
		userState := models.ReturnUserState{
			UserState: states[i].UserState,
			StateID:   states[i].MerklePath.Path,
		}
		userStates = append(userStates, userState)
	}

	return userStates, nil
}
