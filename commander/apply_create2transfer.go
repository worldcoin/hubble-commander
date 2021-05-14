package commander

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyCreate2Transfer(
	storage *st.Storage,
	client *eth.Client,
	create2Transfer *models.Create2Transfer,
	commitmentTokenIndex models.Uint256,
) (addedPubKeyID *uint32, create2TransferError, appError error) {
	pubKeyID, err := getPubKeyID(storage, client, create2Transfer, commitmentTokenIndex)
	if err != nil {
		return nil, nil, err
	}

	stateTree := st.NewStateTree(storage)
	emptyUserState := models.UserState{
		PubKeyID:   *pubKeyID,
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

	return pubKeyID, nil, nil
}

func getPubKeyID(storage *st.Storage, client *eth.Client, transfer *models.Create2Transfer, tokenIndex models.Uint256) (*uint32, error) {
	pubKeyID, err := storage.GetUnusedPubKeyID(&transfer.ToPublicKey, tokenIndex)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}
	if st.IsNotFoundError(err) {
		return client.RegisterAccount(transfer.ToPublicKey)
	}
	return pubKeyID, nil
}
