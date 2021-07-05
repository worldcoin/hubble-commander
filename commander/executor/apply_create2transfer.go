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
	emptyUserState := models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  commitmentTokenID,
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}

	// TODO-AFS isn't ToStateID always nil here?
	if create2Transfer.ToStateID == nil {
		nextAvailableStateID, err := t.storage.GetNextAvailableStateID()
		if err != nil {
			return nil, err
		}
		create2Transfer.ToStateID = nextAvailableStateID
	}

	_, err := t.stateTree.Set(*create2Transfer.ToStateID, &emptyUserState)
	if err != nil {
		return nil, err
	}

	return t.ApplyTransfer(create2Transfer, commitmentTokenID)
}

func (t *TransactionExecutor) ApplyCreate2TransferForSync(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (syncedTransfer *SyncedTransfer, transferError, appError error) {
	emptyUserState := models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  commitmentTokenID,
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}

	// TODO-AFS not necessary
	if create2Transfer.ToStateID == nil {
		nextAvailableStateID, err := t.storage.GetNextAvailableStateID()
		if err != nil {
			return nil, nil, err
		}
		create2Transfer.ToStateID = nextAvailableStateID
	}

	_, err := t.stateTree.Set(*create2Transfer.ToStateID, &emptyUserState)
	if err != nil {
		return nil, nil, err
	}

	return t.ApplyTransferForSync(create2Transfer, commitmentTokenID)
}
