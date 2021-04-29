package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyCreate2Transfer(
	storage *st.Storage,
	create2Transfer *models.Create2Transfer,
	commitmentTokenIndex models.Uint256,
) (addedPubKeyID *uint32, create2TransferError, appError error) {
	stateTree := st.NewStateTree(storage)
	emptyUserState := models.UserState{
		PubKeyID:   create2Transfer.ToPubKeyID,
		TokenIndex: commitmentTokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}

	nextAvailableStateID, err := storage.GetNextAvailableStateID()
	if err != nil {
		return nil, nil, err
	}

	err = stateTree.Set(*nextAvailableStateID, &emptyUserState)
	if err != nil {
		return nil, nil, err
	}

	transfer := models.Transfer{
		TransactionBase: create2Transfer.TransactionBase,
		ToStateID:       *nextAvailableStateID,
	}

	create2TransferError, appError = ApplyTransfer(stateTree, &transfer, commitmentTokenIndex)
	if create2TransferError != nil || appError != nil {
		return nil, create2TransferError, appError
	}

	return &create2Transfer.ToPubKeyID, nil, nil
}
