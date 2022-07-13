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
	// TODO: Rework the interface here. This witness does not include the state root,
	//       and the response does not give any indication as to *which* state root
	//       this proof was built off of.
	//
	//       Also, this method gives results which are inconsistent with
	//       hubble_getUserState, which reads from the pending state, while this
	//       builds a proof off of the batched state.

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
