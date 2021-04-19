package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyCreate2Transfer(
	storage *st.Storage,
	create2transfer *models.Create2Transfer,
	feeReceiverTokenIndex models.Uint256,
) (create2transferError, appError error) {
	stateTree := st.NewStateTree(storage)
	senderUserState, err := stateTree.Leaf(create2transfer.FromStateID)
	if err != nil {
		return nil, err
	}
	emptyUserState := models.UserState{
		PubKeyID:   create2transfer.ToPubKeyID,
		TokenIndex: senderUserState.TokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}

	nextAvailableStateID, err := storage.GetNextAvailableStateID()
	if err != nil {
		return nil, err
	}

	err = stateTree.Set(*nextAvailableStateID, &emptyUserState)
	if err != nil {
		return nil, err
	}

	transfer := models.Transfer{
		TransactionBase: create2transfer.TransactionBase,
		ToStateID:       *nextAvailableStateID,
	}

	create2transferError, appError = ApplyTransfer(stateTree, &transfer, feeReceiverTokenIndex)
	if create2transferError != nil || appError != nil {
		return create2transferError, appError
	}

	return nil, nil
}
