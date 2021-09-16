package applier

import "github.com/Worldcoin/hubble-commander/models"

func (c *Applier) ApplyTransfer(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult SingleTxResult, transferError, appError error,
) {
	receiverLeaf, appError := c.storage.StateTree.Leaf(*tx.GetToStateID())
	if appError != nil {
		return nil, nil, appError
	}

	transferError, appError = c.ApplyTx(tx, receiverLeaf, commitmentTokenID)
	return &ApplySingleTransferResult{Tx: tx}, transferError, appError
}
