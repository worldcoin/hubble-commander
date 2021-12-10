package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionWithBatchDetails struct {
	Transaction interface{}
	BatchHash   *common.Hash
	BatchTime   *models.Timestamp
}

func MakeTransactionWithBatchDetails(tx *models.TransactionWithBatchDetails) TransactionWithBatchDetails {
	return TransactionWithBatchDetails{
		Transaction: tx.Transaction,
		BatchHash:   tx.BatchHash,
		BatchTime:   tx.BatchTime,
	}
}
