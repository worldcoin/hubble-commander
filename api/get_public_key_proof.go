package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getPublicKeyProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(20003, "public key not found"),
}

func (a *API) GetPublicKeyProofByID(id uint32) (*dto.PublicKeyProof, error) {
	publicKeyProof, err := a.unsafeGetPublicKeyProofByID(id)
	if err != nil {
		return nil, sanitizeError(err, getPublicKeyProofAPIErrors)
	}
	return publicKeyProof, nil
}

func (a *API) unsafeGetPublicKeyProofByID(id uint32) (*dto.PublicKeyProof, error) {
	leaf, err := a.storage.AccountTree.Leaf(id)

	if err != nil {
		return nil, err
	}

	witness, err := a.storage.AccountTree.GetWitness(id)

	if err != nil {
		return nil, err
	}

	return &dto.PublicKeyProof{
		PublicKeyProof: models.PublicKeyProof{
			PublicKey: &leaf.PublicKey,
			Witness:   witness,
		},
	}, nil
}
