package applier

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (a *Applier) ApplyTransfer(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult ApplySingleTxResult, transferError, appError error,
) {
	receiverLeaf, appError := a.storage.StateTree.Leaf(*tx.GetToStateID())
	if st.IsNotFoundError(appError) {
		transferError = appError
		return nil, transferError, nil
	}
	if appError != nil {
		return nil, nil, appError
	}

	transferError, appError = a.ApplyTx(tx, receiverLeaf, commitmentTokenID)
	if transferError != nil || appError != nil {
		return nil, transferError, appError
	}
	return &ApplySingleTransferResult{tx: tx.ToTransfer()}, nil, nil
}

func (a *Applier) ApplyTransferForSync(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	synced *SyncedGenericTransaction,
	transferError, appError error,
) {
	receiverLeaf, err := a.storage.StateTree.LeafOrEmpty(*tx.GetToStateID())
	if err != nil {
		return nil, nil, err
	}

	return a.applyTxForSync(tx, receiverLeaf, commitmentTokenID)
}
