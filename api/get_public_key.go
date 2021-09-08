package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getPublicKeyAPIErrors = map[error]ErrorAPI{
	&storage.NotFoundError{}: {
		Code:    10000,                  // TODO-API what here?
		Message: "public key not found", // TODO-API make this more verbose
	},
}

func (a *API) GetPublicKeyByID(id uint32) (*models.PublicKey, error) {
	publicKey, err := a.unsafeGetPublicKeyByID(id)
	if err != nil {
		return nil, sanitizeError(err, getBatchAPIErrors)
	}

	return publicKey, nil
}

func (a *API) unsafeGetPublicKeyByID(id uint32) (*models.PublicKey, error) {
	leaf, err := a.storage.AccountTree.Leaf(id)
	if err != nil {
		return nil, err
	}
	return &leaf.PublicKey, nil
}

func (a *API) GetPublicKeyByStateID(id uint32) (*models.PublicKey, error) {
	publicKey, err := a.storage.GetPublicKeyByStateID(id)
	if err != nil {
		return nil, sanitizeError(err, getBatchAPIErrors)
	}

	return publicKey, nil
}
