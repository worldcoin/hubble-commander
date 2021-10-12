package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetPublicKeyByIDProof(id uint32) (*dto.PublicKeyProof, error) {
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
