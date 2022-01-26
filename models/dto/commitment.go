package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
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

type TxCommitment struct {
	ID                 CommitmentID
	Type               batchtype.BatchType
	PostStateRoot      common.Hash
	LeafHash           common.Hash
	TokenID            models.Uint256
	FeeReceiverStateID uint32
	CombinedSignature  models.Signature
	Status             batchstatus.BatchStatus
	BatchTime          models.Timestamp
	Transactions       interface{}
}

type MMCommitment struct {
	ID                CommitmentID
	Type              batchtype.BatchType
	PostStateRoot     common.Hash
	LeafHash          common.Hash
	CombinedSignature models.Signature
	Status            batchstatus.BatchStatus
	BatchTime         models.Timestamp
	WithdrawRoot      common.Hash
	Meta              MassMigrationMeta
	Transactions      interface{}
}

type DepositCommitment struct {
	ID            CommitmentID
	Type          batchtype.BatchType
	PostStateRoot common.Hash
	LeafHash      common.Hash
	Status        batchstatus.BatchStatus
	BatchTime     models.Timestamp
	SubtreeID     models.Uint256
	SubtreeRoot   common.Hash
	Deposits      []PendingDeposit
}

func NewTxCommitment(
	commitment *models.TxCommitment,
	tokenID models.Uint256,
	status *batchstatus.BatchStatus,
	batchTime *models.Timestamp,
	transactions interface{},
) *TxCommitment {
	return &TxCommitment{
		ID:                 *NewCommitmentID(&commitment.ID),
		Type:               commitment.Type,
		PostStateRoot:      commitment.PostStateRoot,
		LeafHash:           commitment.LeafHash(),
		TokenID:            tokenID,
		FeeReceiverStateID: commitment.FeeReceiver,
		CombinedSignature:  commitment.CombinedSignature,
		Status:             *status,
		BatchTime:          *batchTime,
		Transactions:       transactions,
	}
}

func NewMMCommitment(
	commitment *models.MMCommitment,
	status *batchstatus.BatchStatus,
	batchTime *models.Timestamp,
	transactions interface{},
) *MMCommitment {
	return &MMCommitment{
		ID:                *NewCommitmentID(&commitment.ID),
		Type:              commitment.Type,
		PostStateRoot:     commitment.PostStateRoot,
		LeafHash:          commitment.LeafHash(),
		CombinedSignature: commitment.CombinedSignature,
		Status:            *status,
		BatchTime:         *batchTime,
		WithdrawRoot:      commitment.WithdrawRoot,
		Meta: MassMigrationMeta{
			SpokeID:            commitment.Meta.SpokeID,
			TokenID:            commitment.Meta.TokenID,
			Amount:             commitment.Meta.Amount,
			FeeReceiverStateID: commitment.Meta.FeeReceiver,
		},
		Transactions: transactions,
	}
}

func NewDepositCommitment(
	commitment *models.DepositCommitment,
	status *batchstatus.BatchStatus,
	batchTime *models.Timestamp,
) *DepositCommitment {
	return &DepositCommitment{
		ID:            *NewCommitmentID(&commitment.ID),
		Type:          commitment.Type,
		PostStateRoot: commitment.PostStateRoot,
		LeafHash:      commitment.LeafHash(),
		Status:        *status,
		BatchTime:     *batchTime,
		SubtreeID:     commitment.SubtreeID,
		SubtreeRoot:   commitment.SubtreeRoot,
		Deposits:      MakePendingDeposits(commitment.Deposits),
	}
}
