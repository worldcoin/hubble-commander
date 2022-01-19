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
	ID                 CommitmentID
	Type               batchtype.BatchType
	PostStateRoot      common.Hash
	LeafHash           common.Hash
	TokenID            *models.Uint256   `json:",omitempty"`
	FeeReceiverStateID *uint32           `json:",omitempty"`
	CombinedSignature  *models.Signature `json:",omitempty"`
	Status             txstatus.TransactionStatus
	BatchTime          *models.Timestamp

	massMigrationCommitmentDetails
	depositCommitmentDetails

	Transactions interface{} `json:",omitempty"`
}

type BatchCommitment struct {
	ID                 CommitmentID
	PostStateRoot      common.Hash
	LeafHash           common.Hash
	TokenID            *models.Uint256   `json:",omitempty"`
	FeeReceiverStateID *uint32           `json:",omitempty"`
	CombinedSignature  *models.Signature `json:",omitempty"`

	massMigrationCommitmentDetails
	depositCommitmentDetails
}

type massMigrationCommitmentDetails struct {
	WithdrawRoot *common.Hash       `json:",omitempty"`
	Meta         *MassMigrationMeta `json:",omitempty"`
}

type depositCommitmentDetails struct {
	SubtreeID   *models.Uint256  `json:",omitempty"`
	SubtreeRoot *common.Hash     `json:",omitempty"`
	Deposits    []PendingDeposit `json:",omitempty"`
}

func NewTxCommitment(
	commitment *models.TxCommitment,
	tokenID models.Uint256,
	status *txstatus.TransactionStatus,
	batchTime *models.Timestamp,
	transactions interface{},
) *Commitment {
	return &Commitment{
		ID:                 *NewCommitmentID(&commitment.ID),
		Type:               commitment.Type,
		PostStateRoot:      commitment.PostStateRoot,
		LeafHash:           commitment.LeafHash(),
		TokenID:            &tokenID,
		FeeReceiverStateID: &commitment.FeeReceiver,
		CombinedSignature:  &commitment.CombinedSignature,
		Status:             *status,
		BatchTime:          batchTime,
		Transactions:       transactions,
	}
}

func NewMMCommitment(
	commitment *models.MMCommitment,
	status *txstatus.TransactionStatus,
	batchTime *models.Timestamp,
	transactions interface{},
) *Commitment {
	return &Commitment{
		ID:                *NewCommitmentID(&commitment.ID),
		Type:              commitment.Type,
		PostStateRoot:     commitment.PostStateRoot,
		LeafHash:          commitment.LeafHash(),
		CombinedSignature: &commitment.CombinedSignature,
		Status:            *status,
		BatchTime:         batchTime,
		Transactions:      transactions,
		massMigrationCommitmentDetails: massMigrationCommitmentDetails{
			WithdrawRoot: &commitment.WithdrawRoot,
			Meta: &MassMigrationMeta{
				SpokeID:            commitment.Meta.SpokeID,
				TokenID:            commitment.Meta.TokenID,
				Amount:             commitment.Meta.Amount,
				FeeReceiverStateID: commitment.Meta.FeeReceiver,
			},
		},
	}
}

func NewDepositCommitment(
	commitment *models.DepositCommitment,
	status *txstatus.TransactionStatus,
	batchTime *models.Timestamp,
	transactions interface{},
) *Commitment {
	return &Commitment{
		ID:            *NewCommitmentID(&commitment.ID),
		Type:          commitment.Type,
		PostStateRoot: commitment.PostStateRoot,
		LeafHash:      commitment.LeafHash(),
		Status:        *status,
		BatchTime:     batchTime,
		Transactions:  transactions,
		depositCommitmentDetails: depositCommitmentDetails{
			SubtreeID:   &commitment.SubtreeID,
			SubtreeRoot: &commitment.SubtreeRoot,
			Deposits:    modelsPendingDepositsToDTOPendingDeposits(commitment.Deposits),
		},
	}
}

func MakeTxBatchCommitment(
	commitment *models.TxCommitment,
	tokenID models.Uint256,
) BatchCommitment {
	return BatchCommitment{
		ID:                 *NewCommitmentID(&commitment.ID),
		PostStateRoot:      commitment.PostStateRoot,
		LeafHash:           commitment.LeafHash(),
		TokenID:            &tokenID,
		FeeReceiverStateID: &commitment.FeeReceiver,
		CombinedSignature:  &commitment.CombinedSignature,
	}
}

func MakeMMBatchCommitment(
	commitment *models.MMCommitment,
) BatchCommitment {
	return BatchCommitment{
		ID:                *NewCommitmentID(&commitment.ID),
		PostStateRoot:     commitment.PostStateRoot,
		LeafHash:          commitment.LeafHash(),
		CombinedSignature: &commitment.CombinedSignature,
		massMigrationCommitmentDetails: massMigrationCommitmentDetails{
			WithdrawRoot: &commitment.WithdrawRoot,
			Meta: &MassMigrationMeta{
				SpokeID:            commitment.Meta.SpokeID,
				TokenID:            commitment.Meta.TokenID,
				Amount:             commitment.Meta.Amount,
				FeeReceiverStateID: commitment.Meta.FeeReceiver,
			},
		},
	}
}

func MakeDepositBatchCommitment(
	commitment *models.DepositCommitment,
) BatchCommitment {
	return BatchCommitment{
		ID:            *NewCommitmentID(&commitment.ID),
		PostStateRoot: commitment.PostStateRoot,
		LeafHash:      commitment.LeafHash(),
		depositCommitmentDetails: depositCommitmentDetails{
			SubtreeID:   &commitment.SubtreeID,
			SubtreeRoot: &commitment.SubtreeRoot,
			Deposits:    modelsPendingDepositsToDTOPendingDeposits(commitment.Deposits),
		},
	}
}
