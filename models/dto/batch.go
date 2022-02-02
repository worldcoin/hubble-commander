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
	MinedBlock        *uint32
	SubmissionTime    *models.Timestamp
	Status            batchstatus.BatchStatus
	FinalisationBlock *uint32
}

type BatchWithRootAndCommitments struct {
	Batch
	AccountTreeRoot *common.Hash
	Commitments     interface{}
}

func NewSubmittedBatch(batch *models.Batch) *Batch {
	return &Batch{
		ID:              batch.ID,
		Type:            batch.Type,
		TransactionHash: batch.TransactionHash,
		Status:          batchstatus.Submitted,
	}
}

func NewBatch(batch *models.Batch, minedBlock *uint32, status *batchstatus.BatchStatus) *Batch {
	return &Batch{
		ID:                batch.ID,
		Hash:              batch.Hash,
		Type:              batch.Type,
		TransactionHash:   batch.TransactionHash,
		MinedBlock:        minedBlock,
		SubmissionTime:    batch.SubmissionTime,
		Status:            *status,
		FinalisationBlock: batch.FinalisationBlock,
	}
}

func MakeBatchWithRootAndCommitments(
	batch *Batch,
	accountRoot *common.Hash,
	commitments interface{},
) *BatchWithRootAndCommitments {
	batchDTO := &BatchWithRootAndCommitments{
		Batch:       *batch,
		Commitments: commitments,
	}

	// AccountRoot is always a zero hash for genesis and deposit batches, so we set it to nil
	if batch.Type != batchtype.Genesis && batch.Type != batchtype.Deposit {
		batchDTO.AccountTreeRoot = accountRoot
	}

	return batchDTO
}
