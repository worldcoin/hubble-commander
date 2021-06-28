package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetUserState(id uint32) (*dto.UserState, error) {
	leaf, err := a.storage.GetStateLeaf(id)
	if err != nil {
		return nil, err
	}

	userState := &models.UserStateWithID{
		StateID: leaf.StateID,
		UserState: models.UserState{
			PubKeyID: leaf.PubKeyID,
			TokenID:  leaf.TokenID,
			Balance:  leaf.Balance,
			Nonce:    leaf.Nonce,
		},
	}

	return userState, nil
}
