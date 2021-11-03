package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/ethereum/go-ethereum/common"
)

type Commitment struct {
	models.TxCommitment
	Status       txstatus.TransactionStatus
	BatchTime    *models.Timestamp
	Transactions interface{}
}

type CommitmentWithTokenID struct {
	ID                 models.CommitmentID
	LeafHash           common.Hash
	TokenID            models.Uint256
	FeeReceiverStateID uint32
	CombinedSignature  models.Signature
	PostStateRoot      common.Hash
}

func MakeCommitmentWithTokenID(commitment *models.TxCommitment, tokenID models.Uint256) CommitmentWithTokenID {
	return CommitmentWithTokenID{
		ID:                 commitment.ID,
		LeafHash:           commitment.LeafHash(),
		TokenID:            tokenID,
		FeeReceiverStateID: commitment.FeeReceiver,
		CombinedSignature:  commitment.CombinedSignature,
		PostStateRoot:      commitment.PostStateRoot,
	}
}
