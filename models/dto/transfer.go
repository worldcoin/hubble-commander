package dto

import "github.com/Worldcoin/hubble-commander/models"

type Transfer struct {
	FromStateID *uint32
	ToStateID   *uint32
	Amount      *models.Uint256
	Fee         *models.Uint256
	Nonce       *models.Uint256
	Signature   []byte
}
