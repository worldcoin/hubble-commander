package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

// TODO-AFS consider returning newCreate2Transfer with ToStateID set instead of modifying received parameter
func (t *TransactionExecutor) ApplyCreate2Transfer(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (transferError, appError error) {
	nextAvailableStateID, appError := t.storage.GetNextAvailableStateID()
	if appError != nil {
		return nil, appError
	}
	create2Transfer.ToStateID = nextAvailableStateID

	appError = t.insertNewUserState(*create2Transfer.ToStateID, pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, appError
	}

	return t.ApplyTransfer(create2Transfer, commitmentTokenID)
}

func (t *TransactionExecutor) ApplyCreate2TransferForSync(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (syncedTransfer *SyncedTransfer, transferError, appError error) {
	appError = t.insertNewUserState(*create2Transfer.ToStateID, pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}

	return t.ApplyTransferForSync(create2Transfer, commitmentTokenID)
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
