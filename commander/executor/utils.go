package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
)

func logAndSaveTxError(storage *st.Storage, transaction models.GenericTransaction, txError error) {
	err := storage.SetTransactionError(transaction.GetBase().Hash, txError.Error())
	if err != nil {
		log.Errorf("Setting transaction error failed: %s", err)
	}

	log.Errorf("%s failed: %s", transaction.Type().String(), txError)
}
