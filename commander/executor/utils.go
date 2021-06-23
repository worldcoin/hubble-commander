package executor

import (
	"log"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func logAndSaveTransactionError(storage *st.Storage, transaction *models.TransactionBase, transactionError error) {
	err := storage.SetTransactionError(transaction.Hash, transactionError.Error())
	if err != nil {
		log.Printf("Setting transaction error failed: %s", err)
	}

	log.Printf("%s failed: %s", transaction.TxType.String(), transactionError)
}
