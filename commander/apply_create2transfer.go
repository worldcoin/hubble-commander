package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyCreate2Transfer(
	storage *st.Storage,
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenIndex models.Uint256,
) (create2TransferError, appError error) {
	stateTree := st.NewStateTree(storage)
	emptyUserState := models.UserState{
		PubKeyID:   pubKeyID,
		TokenIndex: commitmentTokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}

	if create2Transfer.ToStateID == nil {
		nextAvailableStateID, err := storage.GetNextAvailableStateID()
		if err != nil {
			return nil, err
		}
		create2Transfer.ToStateID = nextAvailableStateID
	}

	err := stateTree.Set(*create2Transfer.ToStateID, &emptyUserState)
	if err != nil {
		return nil, err
	}

	transfer := models.Transfer{
		TransactionBase: create2Transfer.TransactionBase,
		ToStateID:       *create2Transfer.ToStateID,
	}

	return ApplyTransfer(storage, &transfer, commitmentTokenIndex)
}
