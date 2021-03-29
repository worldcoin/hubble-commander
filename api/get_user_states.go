package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]models.ReturnUserState, error) {
	accounts, err := a.storage.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	userStates := make([]models.ReturnUserState, 0, 1)

	for i := range accounts {
		stateLeaves, err := a.storage.GetStateLeaves(accounts[i].AccountIndex)
		if err != nil {
			return nil, err
		}

		for i := range stateLeaves {
			path, err := a.storage.GetStateNodeByHash(stateLeaves[i].DataHash)
			if err != nil {
				return nil, err
			}
			userState := models.ReturnUserState{
				UserState: stateLeaves[i].UserState,
				StateID:   path.MerklePath.Path,
			}
			userStates = append(userStates, userState)
		}
	}

	return userStates, nil
}
