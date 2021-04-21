package dto

import "github.com/Worldcoin/hubble-commander/models"

type TransferReceipt struct {
	models.Transfer
	Status models.TransactionStatus
}
