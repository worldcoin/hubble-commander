package dto

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
)

type TransactionReceipt struct {
	TransactionWithBatchDetails
	Status txstatus.TransactionStatus
}

type TransferReceipt struct {
	TransferWithBatchDetails
	Status txstatus.TransactionStatus
}

type Create2TransferReceipt struct {
	Create2TransferWithBatchDetails
	Status txstatus.TransactionStatus
}
