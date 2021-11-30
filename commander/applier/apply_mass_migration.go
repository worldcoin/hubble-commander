package applier

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *Applier) ApplyMassMigration(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult ApplySingleTxResult, txError, appError error,
) {
	// Empty receiver leaf used just to pass validation
	// We don't care about the receiver in our state tree for Mass Migration transactions
	dummyReceiverLeaf := &models.StateLeaf{
		UserState: models.UserState{
			TokenID: commitmentTokenID,
		},
	}

	newSenderState, _, txError, appError := a.validateAndCalculateStateAfterTx(tx, dummyReceiverLeaf, commitmentTokenID)
	if txError != nil || appError != nil {
		return nil, txError, appError
	}

	_, appError = a.storage.StateTree.Set(tx.GetFromStateID(), newSenderState)
	if appError != nil {
		return nil, nil, appError
	}

	return &ApplySingleMassMigrationResult{tx: tx.ToMassMigration()}, nil, nil
}
