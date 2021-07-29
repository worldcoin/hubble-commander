package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

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
