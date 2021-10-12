package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetUserStateProof(id uint32) (*dto.StateMerkleProof, error) {
	leaf, err := a.storage.StateTree.Leaf(id)

	if err != nil {
		return nil, err
	}

	witness, err := a.storage.StateTree.GetLeafWitness(id)

	if err != nil {
		return nil, err
	}

	return &dto.StateMerkleProof{
		StateMerkleProof: models.StateMerkleProof{
			UserState: &leaf.UserState,
			Witness:   witness,
		},
	}, nil
}
