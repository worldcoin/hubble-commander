package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyCreate2Transfer(storage *st.Storage, create2transfer *models.Create2Transfer, feeReceiverTokenIndex models.Uint256) (transferError, appError error) {
	stateTree := st.NewStateTree(storage)
	senderUserState, err := stateTree.Leaf(create2transfer.FromStateID)
	if err != nil {
		return nil, err
	}
	emptyUserState := models.UserState{
		PubkeyID:   create2transfer.ToPubkeyID,
		TokenIndex: senderUserState.TokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}

	nextAvailableLeafPath, err := storage.GetNextAvailableLeafPath()
	if err != nil {
		return nil, err
	}

	stateTree.Set(nextAvailableLeafPath.Path, &emptyUserState)

	transfer := models.Transfer{
		TransactionBase: create2transfer.TransactionBase,
		ToStateID:       nextAvailableLeafPath.Path,
	}

	transferError, appError = ApplyTransfer(stateTree, &transfer, feeReceiverTokenIndex)
	if transferError != nil || appError != nil {
		return transferError, appError
	}

	return nil, nil
}
