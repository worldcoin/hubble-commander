package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]models.UserState, error) {
	accounts, err := a.storage.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	userStates := make([]models.UserState, 0, 1)

	for i := range accounts {
		stateLeafs, err := a.storage.GetStateLeafs(accounts[i].AccountIndex)
		if err != nil {
			return nil, err
		}

		for i := range stateLeafs {
			userStates = append(userStates, stateLeafs[i].UserState)
		}
	}

	return userStates, nil
}
