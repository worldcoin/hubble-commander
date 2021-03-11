package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

func ApplyTransactions(
	stateTree *storage.StateTree,
	transactions []models.Transaction,
	feeReceiverIndex uint32,
) (
	[]models.Transaction,
	error,
) {
	validTxs := make([]models.Transaction, 0, 32)

	for i := range transactions {
		tx := transactions[i]
		txError, appError := ApplyTransfer(stateTree, &tx, feeReceiverIndex)
		if appError != nil {
			return nil, appError
		}
		if txError == nil {
			validTxs = append(validTxs, tx)
		} else {
			log.Printf("Transaction failed: %e", txError)
		}

		if len(validTxs) == 32 {
			break
		}
	}

	return validTxs, nil
}
