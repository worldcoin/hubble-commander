package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                models.Uint256
	Hash              *common.Hash
	Type              batchtype.BatchType
	TransactionHash   common.Hash
	SubmissionBlock   *uint32
	SubmissionTime    *models.Timestamp
	Status            batchstatus.BatchStatus
	FinalisationBlock *uint32
}

type BatchWithRootAndCommitments struct {
	Batch
	AccountTreeRoot *common.Hash
	Commitments     interface{}
}

func MakeBatch(batch *models.Batch, submissionBlock uint32, status *batchstatus.BatchStatus) *Batch {
	return &Batch{
		ID:                batch.ID,
		Hash:              batch.Hash,
		Type:              batch.Type,
		TransactionHash:   batch.TransactionHash,
		SubmissionBlock:   submissionBlock,
		SubmissionTime:    batch.SubmissionTime,
		Status:            *status,
		FinalisationBlock: batch.FinalisationBlock,
	}
}

func MakeBatchWithRootAndCommitments(
	batch *models.Batch,
	submissionBlock *uint32,
	status *batchstatus.BatchStatus,
	commitments interface{},
) *BatchWithRootAndCommitments {
	batchDTO := &BatchWithRootAndCommitments{
		Batch:       *MakeBatch(batch, submissionBlock, status),
		Commitments: commitments,
	}

	// AccountRoot is always a zero hash for genesis and deposit batches, so we set it to nil
	if batch.Type != batchtype.Genesis && batch.Type != batchtype.Deposit {
		batchDTO.AccountTreeRoot = batch.AccountTreeRoot
	}

	return batchDTO
}
