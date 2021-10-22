package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getPublicKeyProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(99005, "public key proof not found"),
}

func (a *API) GetPublicKeyProofByPubKeyID(id uint32) (*dto.PublicKeyProof, error) {
	publicKeyProof, err := a.unsafeGetPublicKeyProofByPubKeyID(id)
	if err != nil {
		return nil, sanitizeError(err, getPublicKeyProofAPIErrors)
	}
	return publicKeyProof, nil
}

func (a *API) unsafeGetPublicKeyProofByPubKeyID(id uint32) (*dto.PublicKeyProof, error) {
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
