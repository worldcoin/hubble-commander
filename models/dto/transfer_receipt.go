package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
)

type TransferReceipt struct {
	models.TransferWithBatchDetails
	Status txstatus.TransactionStatus
}

type Create2TransferReceipt struct {
	models.Create2TransferWithBatchDetails
	Status txstatus.TransactionStatus
}

type MassMigrationReceipt struct {
	models.MassMigrationWithBatchDetails
	Status txstatus.TransactionStatus
}
