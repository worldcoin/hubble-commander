package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                models.Uint256
	Hash              *common.Hash
	Type              batchtype.BatchType
	TransactionHash   common.Hash
	SubmissionBlock   uint32
	SubmissionTime    *models.Timestamp
	FinalisationBlock *uint32
}

type BatchWithRootAndCommitments struct {
	Batch
	AccountTreeRoot *common.Hash
	Commitments     []BatchCommitment
}

func MakeBatch(batch *models.Batch, submissionBlock uint32) *Batch {
	return &Batch{
		ID:                batch.ID,
		Hash:              batch.Hash,
		Type:              batch.Type,
		TransactionHash:   batch.TransactionHash,
		SubmissionBlock:   submissionBlock,
		SubmissionTime:    batch.SubmissionTime,
		FinalisationBlock: batch.FinalisationBlock,
	}
}

func MakeBatchWithRootAndCommitments(
	batch *models.Batch,
	submissionBlock uint32,
	commitments []BatchCommitment,
) *BatchWithRootAndCommitments {
	return &BatchWithRootAndCommitments{
		Batch:           *MakeBatch(batch, submissionBlock),
		AccountTreeRoot: batch.AccountTreeRoot,
		Commitments:     commitments,
	}
}
