package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
)

type Commitment struct {
	models.Commitment
	Status       txstatus.TransactionStatus
	BatchTime    *models.Timestamp
	Transactions interface{}
}

type TransferForCommitment struct {
	*models.TransferForCommitment
	ReceiveTime *models.Timestamp
}

type Create2TransferForCommitment struct {
	*models.Create2TransferForCommitment
	ReceiveTime *models.Timestamp
}
