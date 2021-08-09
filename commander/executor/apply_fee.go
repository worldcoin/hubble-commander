package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
)

var ErrInvalidFeeReceiverTokenID = errors.New("invalid fee receiver token ID")

func (t *TransactionExecutor) ApplyFee(feeReceiverStateID uint32, fee models.Uint256) (*models.StateMerkleProof, error) {
	feeReceiver, err := t.storage.StateTree.Leaf(feeReceiverStateID)
	if err != nil {
		return nil, err
	}
	stateProof := &models.StateMerkleProof{
		UserState: feeReceiver.UserState.Copy(),
	}

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateProof.Witness, err = t.storage.StateTree.Set(feeReceiver.StateID, &feeReceiver.UserState)
	if err != nil {
		return nil, err
	}
	return stateProof, nil
}

// TODO use LeafOrEmpty
func (t *TransactionExecutor) ApplyFeeForSync(feeReceiverStateID uint32, commitmentTokenID, fee *models.Uint256) (
	stateProof *models.StateMerkleProof,
	commitmentError error,
	appError error,
) {
	stateProof, appError = t.ApplyFee(feeReceiverStateID, *fee)
	if appError != nil {
		return nil, nil, appError
	}
	if stateProof.UserState.TokenID != *commitmentTokenID {
		return stateProof, ErrInvalidFeeReceiverTokenID, nil
	}
	return stateProof, nil, nil
}
