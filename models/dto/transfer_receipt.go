package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
)

type TransferReceipt struct {
	models.TransferWithBatchDetails
	ReceiveTime *models.Timestamp
	Status      txstatus.TransactionStatus
}

type Create2TransferReceipt struct {
	models.Create2TransferWithBatchDetails
	ReceiveTime *models.Timestamp
	Status      txstatus.TransactionStatus
}
