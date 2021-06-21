package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) ApplyCreate2Transfer(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenIndex models.Uint256,
) (create2TransferError, appError error) {
	emptyUserState := models.UserState{
		PubKeyID:   pubKeyID,
		TokenIndex: commitmentTokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
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

	return t.ApplyTransfer(create2Transfer, commitmentTokenIndex)
}
