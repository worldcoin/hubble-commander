package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var getUserStateProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(50003, "user state proof not found"),
}

func (a *API) GetUserStateProof(id uint32) (*dto.StateMerkleProof, error) {
	if !a.cfg.EnableProofEndpoints {
		return nil, errProofEndpointsDisabled
	}
	userStateProof, err := a.unsafeGetUserStateProof(id)
	if err != nil {
		return nil, sanitizeError(err, getUserStateProofAPIErrors)
	}
	return userStateProof, nil
}

func (a *API) unsafeGetUserStateProof(id uint32) (*dto.StateMerkleProof, error) {
	leaf, err := a.storage.StateTree.Leaf(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	witness, err := a.storage.StateTree.GetLeafWitness(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &dto.StateMerkleProof{
		StateMerkleProof: models.StateMerkleProof{
			UserState: &leaf.UserState,
			Witness:   witness,
		},
	}, nil
}
