package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (t *TransactionExecutor) ApplyFee(feeReceiverStateID uint32, fee models.Uint256) error {
	feeReceiver, err := t.Storage.GetStateLeaf(feeReceiverStateID)
	if err != nil {
		return err
	}

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateTree := st.NewStateTree(t.Storage)
	_, err = stateTree.Set(feeReceiver.StateID, &feeReceiver.UserState)
	return err
}
