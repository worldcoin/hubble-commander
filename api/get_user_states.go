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
		stateLeafs, err := a.storage.GetStateLeafs(accounts[i].AccountIndex)
		if err != nil {
			return nil, err
		}

		for i := range stateLeafs {
			path, err := a.storage.GetStateNodeByHash(stateLeafs[i].DataHash)
			if err != nil {
				return nil, err
			}
			userState := models.ReturnUserState{
				UserState: stateLeafs[i].UserState,
				StateID:   path.MerklePath.Path,
			}
			userStates = append(userStates, userState)
		}
	}

	return userStates, nil
}
