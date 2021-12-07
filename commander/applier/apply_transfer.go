package applier

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var ErrNonexistentReceiver = fmt.Errorf("receiver state ID does not exist")

func (a *Applier) ApplyTransfer(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult ApplySingleTxResult, txError, appError error,
) {
	receiverLeaf, appError := a.storage.StateTree.Leaf(*tx.GetToStateID())
	if st.IsNotFoundError(appError) {
		txError = errors.WithStack(ErrNonexistentReceiver)
		return nil, txError, nil
	}
	if appError != nil {
		return nil, nil, appError
	}

	txError, appError = a.ApplyTx(tx, receiverLeaf, commitmentTokenID)
	if txError != nil || appError != nil {
		return nil, txError, appError
	}
	return &ApplySingleTransferResult{tx: tx.ToTransfer()}, nil, nil
}

func (a *Applier) ApplyTransferForSync(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	synced *SyncedTxWithProofs,
	txError, appError error,
) {
	receiverLeaf, err := a.storage.StateTree.LeafOrEmpty(*tx.GetToStateID())
	if err != nil {
		return nil, nil, err
	}

	return a.applyTxForSync(tx, receiverLeaf, commitmentTokenID)
}
