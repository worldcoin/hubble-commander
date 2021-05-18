package api

import "github.com/Worldcoin/hubble-commander/models"

func (a *API) GetUserState(id uint32) (*models.UserStateWithID, error) {
	leaf, err := a.storage.GetStateLeafByStateID(id)
	if err != nil {
		return nil, err
	}

	userState := &models.UserStateWithID{
		StateID: leaf.StateID,
		UserState: models.UserState{
			PubKeyID:   leaf.PubKeyID,
			TokenIndex: leaf.TokenIndex,
			Balance:    leaf.Balance,
			Nonce:      leaf.Nonce,
		},
	}

	return userState, nil
}
