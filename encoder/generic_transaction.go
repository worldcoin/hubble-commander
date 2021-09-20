package encoder

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
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

func GetTransactionLength(batchType batchtype.BatchType) int {
	switch batchType {
	case batchtype.Transfer:
		return TransferLength
	case batchtype.Create2Transfer:
		return Create2TransferLength
	case batchtype.Genesis, batchtype.MassMigration, batchtype.Deposit:
		log.Panicf("unsupported tx type: %s", batchType)
	}
	return -1
}
