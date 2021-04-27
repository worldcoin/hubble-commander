package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func markTransactionsAsIncluded(storage *st.Storage, transactions []models.TransactionBase, commitmentID int32) error {
	for i := range transactions {
		err := storage.MarkTransactionAsIncluded(transactions[i].Hash, commitmentID)
		if err != nil {
			return err
		}
	}
	return nil
}
