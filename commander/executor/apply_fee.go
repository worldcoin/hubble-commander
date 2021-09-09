package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
)

var ErrInvalidFeeReceiverTokenID = errors.New("invalid fee receiver token ID")

func (t *ExecutionContext) ApplyFee(feeReceiverStateID uint32, fee models.Uint256) (*models.StateMerkleProof, error) {
	feeReceiver, err := t.storage.StateTree.Leaf(feeReceiverStateID)
	if err != nil {
		return nil, err
	}
	return t.applyFee(feeReceiver, fee)
}

func (t *ExecutionContext) applyFee(feeReceiver *models.StateLeaf, fee models.Uint256) (*models.StateMerkleProof, error) {
	initialState := feeReceiver.UserState.Copy()

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateChangeWitness, err := t.storage.StateTree.Set(feeReceiver.StateID, &feeReceiver.UserState)
	if err != nil {
		return nil, err
	}

	stateProof := &models.StateMerkleProof{
		UserState: initialState,
		Witness:   stateChangeWitness,
	}
	return stateProof, nil
}

func (t *ExecutionContext) ApplyFeeForSync(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof,
	commitmentError error,
	appError error,
) {
	feeReceiver, appError := t.storage.StateTree.LeafOrEmpty(feeReceiverStateID)
	if appError != nil {
		return nil, nil, appError
	}
	stateProof, appError = t.applyFee(feeReceiver, *fee)
	if appError != nil {
		return nil, nil, appError
	}

	if stateProof.UserState.TokenID != *commitmentTokenID {
		return stateProof, ErrInvalidFeeReceiverTokenID, nil
	}

	return stateProof, nil, nil
}
