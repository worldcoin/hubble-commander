package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyCreate2Transfer(
	storage *st.Storage,
	create2Transfer *models.Create2Transfer,
	feeReceiverTokenIndex models.Uint256,
) (create2TransferError, appError error) {
	stateTree := st.NewStateTree(storage)
	senderUserState, err := stateTree.Leaf(create2Transfer.FromStateID)
	if err != nil {
		return nil, err
	}
	emptyUserState := models.UserState{
		PubkeyID:   create2Transfer.ToPubkeyID,
		TokenIndex: senderUserState.TokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}

	nextAvailableLeafPath, err := storage.GetNextAvailableLeafPath()
	if err != nil {
		return nil, err
	}

	err = stateTree.Set(*nextAvailableLeafPath, &emptyUserState)
	if err != nil {
		return nil, err
	}

	transfer := models.Transfer{
		TransactionBase: create2Transfer.TransactionBase,
		ToStateID:       *nextAvailableLeafPath,
	}

	create2TransferError, appError = ApplyTransfer(stateTree, &transfer, feeReceiverTokenIndex)
	if create2TransferError != nil || appError != nil {
		return create2TransferError, appError
	}

	return nil, nil
}
