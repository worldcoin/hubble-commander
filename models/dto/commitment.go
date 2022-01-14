package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentID struct {
	BatchID      models.Uint256
	IndexInBatch uint8
}

func NewCommitmentID(id *models.CommitmentID) *CommitmentID {
	if id == nil {
		return nil
	}
	return &CommitmentID{
		BatchID:      id.BatchID,
		IndexInBatch: id.IndexInBatch,
	}
}

type Commitment struct {
	ID                CommitmentID
	Type              batchtype.BatchType
	PostStateRoot     common.Hash
	FeeReceiver       uint32
	CombinedSignature models.Signature
	Status            txstatus.TransactionStatus
	BatchTime         *models.Timestamp
	Transactions      interface{}
}

type CommitmentWithTokenID struct {
	ID                 CommitmentID
	LeafHash           common.Hash
	TokenID            models.Uint256
	FeeReceiverStateID uint32
	CombinedSignature  models.Signature
	PostStateRoot      common.Hash
	Meta               *MassMigrationMeta `json:",omitempty"`
	WithdrawRoot       *common.Hash       `json:",omitempty"`
}

func MakeCommitmentWithTokenIDFromTxCommitment(commitment *models.TxCommitment, tokenID models.Uint256) CommitmentWithTokenID {
	return CommitmentWithTokenID{
		ID:                 *NewCommitmentID(&commitment.ID),
		LeafHash:           commitment.LeafHash(),
		TokenID:            tokenID,
		FeeReceiverStateID: commitment.FeeReceiver,
		CombinedSignature:  commitment.CombinedSignature,
		PostStateRoot:      commitment.PostStateRoot,
		Meta:               nil,
		WithdrawRoot:       nil,
	}
}

func MakeCommitmentWithTokenIDFromMMCommitment(commitment *models.MMCommitment, tokenID models.Uint256) CommitmentWithTokenID {
	return CommitmentWithTokenID{
		ID:                 *NewCommitmentID(&commitment.ID),
		LeafHash:           commitment.LeafHash(),
		TokenID:            tokenID,
		FeeReceiverStateID: commitment.FeeReceiver,
		CombinedSignature:  commitment.CombinedSignature,
		PostStateRoot:      commitment.PostStateRoot,
		Meta:               NewMassMigrationMeta(commitment.Meta),
		WithdrawRoot:       &commitment.WithdrawRoot,
	}
}
