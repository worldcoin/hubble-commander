package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetUserStates(publicKey *models.PublicKey) ([]dto.UserState, error) {
	return a.storage.GetUserStatesByPublicKey(publicKey)
}
