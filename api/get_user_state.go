package api

import "github.com/Worldcoin/hubble-commander/models"

func (a *API) GetUserState(id uint32) (*models.UserStateWithID, error) {
	return a.storage.GetUserStateByID(id)
}
