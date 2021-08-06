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
	userState := feeReceiver.UserState
	stateProof := &models.StateMerkleProof{UserState: &userState}

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateProof.Witness, err = t.storage.StateTree.Set(feeReceiver.StateID, &feeReceiver.UserState)
	if err != nil {
		return nil, err
	}
	return stateProof, err
}

func (t *TransactionExecutor) ApplyFeeForSync(
	feeReceiver *FeeReceiver,
	tokenID, fee *models.Uint256,
) (*models.StateMerkleProof, error, error) {
	stateProof, err := t.ApplyFee(feeReceiver.StateID, *fee)
	if err != nil {
		return nil, nil, err
	}
	if feeReceiver.TokenID != *tokenID {
		return stateProof, ErrInvalidFeeReceiverTokenID, nil
	}
	return stateProof, nil, nil
}
