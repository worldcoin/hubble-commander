package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type ResetPubkey struct {
	NewAccountTreeRoot common.Hash
	OldPubKey          *models.PublicKey
}
