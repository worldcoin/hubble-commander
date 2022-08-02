package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var getUserStateProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(50003, "user state inclusion proof could not be generated"),
}

func (a *API) GetUserStateProof(id uint32) (*dto.StateMerkleProof, error) {
	// TODO: I'm not sure how anybody is supposed to use this since it does not tell
	//       You which root you're supposed to check the proof against. Soon it will
	//       give results which are inconsistent with the results hubble_getUserState
	//       returns. I should check whether that inconsistency breaks any use-cases

	if !a.cfg.EnableProofMethods {
		return nil, APIErrProofMethodsDisabled
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

	dtoUserState := dto.MakeUserState(&leaf.UserState)

	return &dto.StateMerkleProof{
		UserState: &dtoUserState,
		Witness:   witness,
	}, nil
}
