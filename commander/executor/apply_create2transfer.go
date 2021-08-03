package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) ApplyCreate2Transfer(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (appliedTransfer *models.Create2Transfer, transferError, appError error) {
	nextAvailableStateID, appError := t.storage.StateTree.NextAvailableStateID()
	if appError != nil {
		return appliedTransfer, nil, appError
	}
	appliedTransfer = create2Transfer.Clone()
	appliedTransfer.ToStateID = nextAvailableStateID

	receiverLeaf := newUserLeaf(*appliedTransfer.ToStateID, pubKeyID, commitmentTokenID)
	transferError, appError = t.ApplyTransfer(appliedTransfer, receiverLeaf, commitmentTokenID)
	return appliedTransfer, transferError, appError
}

func (t *TransactionExecutor) ApplyCreate2TransferForSync(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (synced *SyncedCreate2Transfer, transferError, appError error) {
	if create2Transfer.ToStateID == nil {
		return nil, nil, ErrNilReceiverStateID
	}

	receiverLeaf := newUserLeaf(*create2Transfer.ToStateID, pubKeyID, commitmentTokenID)
	genericSynced, transferError, appError := t.applyGenericTransactionForSync(create2Transfer, receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	return NewSyncedCreate2TransferFromGeneric(genericSynced), transferError, nil
}

func newUserLeaf(stateID, pubKeyID uint32, tokenID models.Uint256) *models.StateLeaf {
	return &models.StateLeaf{
		StateID: stateID,
		UserState: models.UserState{
			PubKeyID: pubKeyID,
			TokenID:  tokenID,
			Balance:  models.MakeUint256(0),
			Nonce:    models.MakeUint256(0),
		},
	}
}
