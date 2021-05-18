package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func logAndSaveTransactionError(storage *st.Storage, transaction *models.TransactionBase, transactionError error) {
	if transactionError != nil {
		err := storage.SetTransactionError(transaction.Hash, transactionError.Error())
		if err != nil {
			log.Printf("Setting transaction error failed: %s", err)
		}

		log.Printf("%s failed: %s", transaction.TxType.String(), transactionError)
	}
}

func ApplyFee(storage *st.Storage, feeReceiverPubKeyID uint32, tokenIndex, fee models.Uint256) (*uint32, error) {
	feeReceiver, err := storage.GetUserStateByPubKeyIDAndTokenIndex(feeReceiverPubKeyID, tokenIndex)
	if err != nil {
		return nil, err
	}

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	stateTree := st.NewStateTree(storage)
	if err := stateTree.Set(feeReceiver.StateID, &feeReceiver.UserState); err != nil {
		return nil, err
	}

	return &feeReceiver.StateID, nil
}
