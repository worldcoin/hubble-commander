package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
)

type TransferReceipt struct {
	models.TransferWithBatchHash
	ReceiveTime *Timestamp
	Status      txstatus.TransactionStatus
}

type Create2TransferReceipt struct {
	models.Create2TransferWithBatchHash
	ReceiveTime *Timestamp
	Status      txstatus.TransactionStatus
}
