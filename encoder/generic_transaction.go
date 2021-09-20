package encoder

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func HashGenericTransaction(tx models.GenericTransaction) (*common.Hash, error) {
	switch x := tx.(type) {
	case *models.Transfer:
		return HashTransfer(x)
	case *models.Create2Transfer:
		return HashCreate2Transfer(x)
	default:
		return nil, errors.Errorf("unsupported tx type: %s", tx.Type())
	}
}

func GetTransactionLength(txType txtype.TransactionType) int {
	switch txType {
	case txtype.Transfer:
		return TransferLength
	case txtype.Create2Transfer:
		return Create2TransferLength
	case txtype.Genesis, txtype.MassMigration, txtype.Deposit:
		log.Panicf("unsupported tx type: %s", txType)
	}
	return -1
}
