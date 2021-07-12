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
