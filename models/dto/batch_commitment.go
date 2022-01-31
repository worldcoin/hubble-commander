package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type BatchTxCommitment struct {
	ID                 CommitmentID
	PostStateRoot      common.Hash
	LeafHash           *common.Hash
	TokenID            models.Uint256
	FeeReceiverStateID uint32
	CombinedSignature  models.Signature
}

type BatchMMCommitment struct {
	ID                CommitmentID
	PostStateRoot     common.Hash
	LeafHash          *common.Hash
	CombinedSignature models.Signature
	WithdrawRoot      common.Hash
	Meta              MassMigrationMeta
}

type BatchDepositCommitment struct {
	ID            CommitmentID
	PostStateRoot common.Hash
	LeafHash      *common.Hash
	SubtreeID     models.Uint256
	SubtreeRoot   common.Hash
	Deposits      []PendingDeposit
}

func MakeBatchTxCommitment(
	commitment *models.TxCommitment,
	tokenID models.Uint256,
) BatchTxCommitment {
	return BatchTxCommitment{
		ID:                 *NewCommitmentID(&commitment.ID),
		PostStateRoot:      commitment.PostStateRoot,
		LeafHash:           LeafHashOrNil(commitment),
		TokenID:            tokenID,
		FeeReceiverStateID: commitment.FeeReceiver,
		CombinedSignature:  commitment.CombinedSignature,
	}
}

func MakeBatchMMCommitment(
	commitment *models.MMCommitment,
) BatchMMCommitment {
	return BatchMMCommitment{
		ID:                *NewCommitmentID(&commitment.ID),
		PostStateRoot:     commitment.PostStateRoot,
		LeafHash:          LeafHashOrNil(commitment),
		CombinedSignature: commitment.CombinedSignature,
		WithdrawRoot:      commitment.WithdrawRoot,
		Meta: MassMigrationMeta{
			SpokeID:            commitment.Meta.SpokeID,
			TokenID:            commitment.Meta.TokenID,
			Amount:             commitment.Meta.Amount,
			FeeReceiverStateID: commitment.Meta.FeeReceiver,
		},
	}
}

func MakeBatchDepositCommitment(
	commitment *models.DepositCommitment,
) BatchDepositCommitment {
	return BatchDepositCommitment{
		ID:            *NewCommitmentID(&commitment.ID),
		PostStateRoot: commitment.PostStateRoot,
		LeafHash:      LeafHashOrNil(commitment),
		SubtreeID:     commitment.SubtreeID,
		SubtreeRoot:   commitment.SubtreeRoot,
		Deposits:      MakePendingDeposits(commitment.Deposits),
	}
}
