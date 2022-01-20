package dto

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
)

type TransactionReceipt struct {
	TransactionWithBatchDetails
	Status txstatus.TransactionStatus
}
