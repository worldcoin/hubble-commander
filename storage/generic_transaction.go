package storage

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
)

func (s *TransactionStorage) BatchAddGenericTransaction(txs models.GenericTransactionArray) error {
	switch x := txs.(type) {
	case models.TransferArray:
		return s.BatchAddTransfer(x)
	case models.Create2TransferArray:
		return s.BatchAddCreate2Transfer(x)
	default:
		return fmt.Errorf("tx type: %t", x) // TODO-API extract the fuck?
	}
}
