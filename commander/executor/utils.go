package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
)

func logAndSaveTransactionError(storage *st.InternalStorage, transaction *models.TransactionBase, transactionError error) {
	err := storage.SetTransactionError(transaction.Hash, transactionError.Error())
	if err != nil {
		log.Printf("Setting transaction error failed: %s", err)
	}

	log.Printf("%s failed: %s", transaction.TxType.String(), transactionError)
}
