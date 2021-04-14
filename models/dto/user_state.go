package dto

import "github.com/Worldcoin/hubble-commander/models"

type ReturnUserState struct {
	StateID uint32
	models.UserState
}
