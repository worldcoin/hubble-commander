package applier

import "github.com/Worldcoin/hubble-commander/models"

func (a *Applier) ApplyTransfer(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult SingleTxResult, transferError, appError error,
) {
	receiverLeaf, appError := a.storage.StateTree.Leaf(*tx.GetToStateID())
	if appError != nil {
		return nil, nil, appError
	}

	transferError, appError = a.ApplyTx(tx, receiverLeaf, commitmentTokenID)
	return &ApplySingleTransferResult{tx: tx}, transferError, appError
}
