package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (t *TransactionExecutor) ApplyFee(feeReceiverStateID uint32, fee models.Uint256) (*models.StateMerkleProof, error) {
	feeReceiver, err := t.storage.GetStateLeaf(feeReceiverStateID)
	if err != nil {
		return nil, err
	}
	userState := feeReceiver.UserState
	stateProof := &models.StateMerkleProof{UserState: &userState}

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateTree := st.NewStateTree(t.storage)
	stateProof.Witness, err = stateTree.Set(feeReceiver.StateID, &feeReceiver.UserState)
	if err != nil {
		return nil, err
	}
	return stateProof, err
}
