package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetUserStates(publicKey models.PublicKey) ([]models.UserState, error) {
	accounts, err := a.storage.GetAccounts(&publicKey)
	if err != nil {
		return nil, err
	}

	userStates := make([]models.UserState, 0, 1)

	for _, account := range accounts {
		stateLeafs, err := a.storage.GetStateLeafs(account.AccountIndex)
		if err != nil {
			return nil, err
		}

		for _, stateLeaf := range stateLeafs {
			userStates = append(userStates, stateLeaf.UserState)

		}
	}

	return userStates, nil
}
