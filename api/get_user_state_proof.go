package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getUserStateProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(20002, "user state not found"),
}

func (a *API) GetUserStateProof(id uint32) (*dto.StateMerkleProof, error) {
	userStateProof, err := a.unsafeGetUserStateProof(id)
	if err != nil {
		return nil, sanitizeError(err, getUserStateProofAPIErrors)
	}
	return userStateProof, nil
}

func (a *API) unsafeGetUserStateProof(id uint32) (*dto.StateMerkleProof, error) {
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