package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionWithBatchDetails struct {
	TransactionBase
	ToStateID   *uint32           `json:",omitempty"`
	ToPublicKey *models.PublicKey `json:",omitempty"`
	SpokeID     *uint32           `json:",omitempty"`
	BatchHash   *common.Hash
	BatchTime   *models.Timestamp
}

func MakeTransactionWithBatchDetails(tx *models.TransactionWithBatchDetails) TransactionWithBatchDetails {
	dtoTx := TransactionWithBatchDetails{
		BatchHash: tx.BatchHash,
		BatchTime: tx.BatchTime,
	}
	switch subTx := tx.Transaction.(type) {
	case *models.Transfer:
		dtoTx.TransactionBase = MakeTransactionBase(&subTx.TransactionBase)
		dtoTx.ToStateID = &subTx.ToStateID
	case *models.Create2Transfer:
		dtoTx.TransactionBase = MakeTransactionBase(&subTx.TransactionBase)
		dtoTx.ToStateID = subTx.ToStateID
		dtoTx.ToPublicKey = &subTx.ToPublicKey
	case *models.MassMigration:
		dtoTx.TransactionBase = MakeTransactionBase(&subTx.TransactionBase)
		dtoTx.SpokeID = &subTx.SpokeID
	}
	return dtoTx
}
