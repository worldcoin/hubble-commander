package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyTransactions(
	storage *st.Storage,
	transactions []models.Transaction,
	feeReceiverIndex uint32,
) (
	[]models.Transaction,
	error,
) {
	stateTree := st.NewStateTree(storage)
	validTxs := make([]models.Transaction, 0, 32)
	combinedFee := models.MakeUint256(0)

	for i := range transactions {
		tx := transactions[i]
		txError, appError := ApplyTransfer(stateTree, &tx)
		if appError != nil {
			return nil, appError
		}
		if txError == nil {
			validTxs = append(validTxs, tx)
			combinedFee.Add(&combinedFee.Int, &tx.Fee.Int)
		} else {
			storage.SetTransactionError(tx.Hash, txError.Error())
			log.Printf("Transaction failed: %s", txError)
		}

		if len(validTxs) == 32 {
			break
		}
	}

	err := ApplyFee(stateTree, feeReceiverIndex, combinedFee)
	if err != nil {
		return nil, err
	}

	return validTxs, nil
}
