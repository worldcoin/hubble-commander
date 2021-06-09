package api

import "github.com/Worldcoin/hubble-commander/models"

func (a *API) GetPublicKeyByID(id uint32) (*models.PublicKey, error) {
	return a.storage.GetPublicKey(id)
}

func (a *API) GetPublicKeyByStateID(id uint32) (*models.PublicKey, error) {
	return a.storage.GetPublicKeyByStateID(id)
}
