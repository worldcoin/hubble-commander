package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
)

type TransferReceipt struct {
	models.Transfer
	Status txstatus.TransactionStatus
}
