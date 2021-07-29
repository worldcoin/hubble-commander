package api

import "github.com/Worldcoin/hubble-commander/models"

func (a *API) GetPublicKeyByID(id uint32) (*models.PublicKey, error) {
	leaf, err := a.storage.AccountTree.Leaf(id)
	if err != nil {
		return nil, err
	}
	return &leaf.PublicKey, nil
}

func (a *API) GetPublicKeyByStateID(id uint32) (*models.PublicKey, error) {
	return a.storage.GetPublicKeyByStateID(id)
}
