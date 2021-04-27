package dto

import "github.com/Worldcoin/hubble-commander/models"

type Commitment struct {
	models.Commitment
	Transactions interface{}
}
