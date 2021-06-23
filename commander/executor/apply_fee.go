package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (t *TransactionExecutor) ApplyFee(tokenIndex, fee models.Uint256) (*uint32, error) {
	feeReceiver, err := t.storage.GetFeeReceiverStateLeaf(t.cfg.FeeReceiverPubKeyID, tokenIndex)
	if err != nil {
		return nil, err
	}

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateTree := st.NewStateTree(t.storage)
	if err := stateTree.Set(feeReceiver.StateID, &feeReceiver.UserState); err != nil {
		return nil, err
	}

	return &feeReceiver.StateID, nil
}
