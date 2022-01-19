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
	batchDTO := &BatchWithRootAndCommitments{
		Batch:       *MakeBatch(batch, submissionBlock),
		Commitments: commitments,
	}

	// AccountRoot is always a zero hash for genesis and deposit batches, so we set it to nil
	if batch.Type != batchtype.Genesis && batch.Type != batchtype.Deposit {
		batchDTO.AccountTreeRoot = batch.AccountTreeRoot
	}

	return batchDTO
}
