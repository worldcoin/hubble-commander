package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetUserStateProof(id uint32) (*models.StateMerkleProof, error) {
	leaf, err := a.storage.StateTree.Leaf(id)

	if err != nil {
		return nil, err
	}

	witness, err := a.storage.StateTree.GetWitness(id)

	if err != nil {
		return nil, err
	}

	return &models.StateMerkleProof{
		UserState: &leaf.UserState,
		Witness:   witness,
	}, nil
}
