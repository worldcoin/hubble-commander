package applier

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
)

var ErrInvalidFeeReceiverTokenID = fmt.Errorf("invalid fee receiver token ID")

func (a *Applier) ApplyFee(feeReceiverStateID uint32, fee models.Uint256) (*models.StateMerkleProof, error) {
	feeReceiver, err := a.storage.StateTree.Leaf(feeReceiverStateID)
	if err != nil {
		return nil, err
	}
	return a.applyFee(feeReceiver, fee)
}

func (a *Applier) applyFee(feeReceiver *models.StateLeaf, fee models.Uint256) (*models.StateMerkleProof, error) {
	initialState := feeReceiver.UserState.Copy()

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateChangeWitness, err := a.storage.StateTree.Set(feeReceiver.StateID, &feeReceiver.UserState)
	if err != nil {
		return nil, err
	}

	stateProof := &models.StateMerkleProof{
		UserState: initialState,
		Witness:   stateChangeWitness,
	}
	return stateProof, nil
}

func (a *Applier) ApplyFeeForSync(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof,
	commitmentError error,
	appError error,
) {
	feeReceiver, appError := a.storage.StateTree.LeafOrEmpty(feeReceiverStateID)
	if appError != nil {
		return nil, nil, appError
	}
	stateProof, appError = a.applyFee(feeReceiver, *fee)
	if appError != nil {
		return nil, nil, appError
	}

	if stateProof.UserState.TokenID != *commitmentTokenID {
		return stateProof, errors.WithStack(ErrInvalidFeeReceiverTokenID), nil
	}

	return stateProof, nil, nil
}
