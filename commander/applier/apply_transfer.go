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

func (a *Applier) ApplyTransferForSync(transfer *models.Transfer, commitmentTokenID models.Uint256) (
	synced *SyncedTransfer,
	transferError, appError error,
) {
	receiverLeaf, err := a.storage.StateTree.LeafOrEmpty(*transfer.GetToStateID())
	if err != nil {
		return nil, nil, err
	}

	genericSynced, transferError, appError := a.applyTxForSync(transfer, receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	return NewSyncedTransferFromGeneric(genericSynced), transferError, nil
}
