package storage

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (s *Storage) BatchAddGenericTransaction(txs models.GenericTransactionArray) error {
	switch txs.Type() {
	case txtype.Transfer:
		return s.BatchAddTransfer(txs.ToTransferArray())
	case txtype.Create2Transfer:
		return s.BatchAddTransfer(txs.ToTransferArray())
	default:
		return fmt.Errorf("unsupported batch type for sync: %s", txs.Type())
	}
}
