package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

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

	if create2Transfer.ToStateID == nil {
		nextAvailableStateID, err := t.storage.GetNextAvailableStateID()
		if err != nil {
			return nil, err
		}
		create2Transfer.ToStateID = nextAvailableStateID
	}

	err := t.stateTree.Set(*create2Transfer.ToStateID, &emptyUserState)
	if err != nil {
		return nil, err
	}

	if t.opts.AssumeNonces { // TODO-AFS rework this
		_, transferError, appError = t.ApplyTransferForSync(create2Transfer, commitmentTokenID)
		return transferError, appError
	} else {
		return t.ApplyTransfer(create2Transfer, commitmentTokenID)
	}
}
