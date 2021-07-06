package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) ApplyCreate2Transfer(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (appliedTransfer *models.Create2Transfer, transferError, appError error) {
	nextAvailableStateID, appError := t.storage.GetNextAvailableStateID()
	if appError != nil {
		return appliedTransfer, nil, appError
	}
	appliedTransfer = create2Transfer.Clone()
	appliedTransfer.ToStateID = nextAvailableStateID

	appError = t.insertNewUserState(*appliedTransfer.ToStateID, pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}

	transferError, appError = t.ApplyTransfer(appliedTransfer, commitmentTokenID)
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

	appError = t.insertNewUserState(*create2Transfer.ToStateID, pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}

	genericSynced, transferError, appError := t.applyGenericTransactionForSync(create2Transfer, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	return NewSyncedCreate2TransferFromGeneric(genericSynced), transferError, nil
}

func (t *TransactionExecutor) insertNewUserState(stateID, pubKeyID uint32, tokenID models.Uint256) error {
	emptyUserState := models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  tokenID,
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}

	_, err := t.stateTree.Set(stateID, &emptyUserState)
	return err
}
