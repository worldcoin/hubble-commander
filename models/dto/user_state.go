package dto

import "github.com/Worldcoin/hubble-commander/models"

type UserStateReceipt struct {
	StateID uint32
	models.UserState
}
