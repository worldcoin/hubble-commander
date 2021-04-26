package commander

import (
	"errors"
	"log"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

var (
	ErrNegativeFee = errors.New("the fee cannot be negative")
)

func logAndSaveTransactionError(storage *st.Storage, transaction *models.TransactionBase, transactionError error) {
	err := storage.SetTransactionError(transaction.Hash, transactionError.Error())
	if err != nil {
		log.Printf("Setting transaction error failed: %s", err)
	}

	log.Printf("%s failed: %s", transaction.TxType.String(), transactionError)
}

func ValidateAndApplyFee(stateTree *st.StateTree, feeReceiverIndex uint32, fee models.Uint256) error {
	if fee.CmpN(0) == 1 {
		// TODO cfg.FeeReceiverIndex actually represents PubKeyID and is used as StateID here
		err := ApplyFee(stateTree, feeReceiverIndex, fee)
		if err != nil {
			return err
		}
	}

	return ErrNegativeFee
}

func ApplyFee(stateTree *st.StateTree, feeReceiverIndex uint32, fee models.Uint256) error {
	feeReceiver, err := stateTree.Leaf(feeReceiverIndex)
	if err != nil {
		return err
	}

	feeReceiver.Balance = *feeReceiver.Balance.Add(&fee)

	err = stateTree.Set(feeReceiverIndex, &feeReceiver.UserState)
	if err != nil {
		return err
	}

	return nil
}
