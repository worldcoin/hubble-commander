package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionBase struct {
	Hash         common.Hash
	TxType       txtype.TransactionType
	FromStateID  uint32
	Amount       models.Uint256
	Fee          models.Uint256
	Nonce        models.Uint256
	Signature    models.Signature
	ReceiveTime  *models.Timestamp
	CommitmentID *CommitmentID
	ErrorMessage *string
}
