package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                models.Uint256
	Hash              *common.Hash
	Type              txtype.TransactionType
	TransactionHash   common.Hash
	SubmissionBlock   uint32
	FinalisationBlock *uint32
}

type BatchWithRootAndCommitments struct {
	Batch
	AccountTreeRoot *common.Hash
	Commitments     []models.CommitmentWithTokenID
}

func MakeBatch(batch *models.Batch, submissionBlock uint32) *Batch {
	return &Batch{
		ID:                batch.Number,
		Hash:              batch.Hash,
		Type:              batch.Type,
		TransactionHash:   batch.TransactionHash,
		SubmissionBlock:   submissionBlock,
		FinalisationBlock: batch.FinalisationBlock,
	}
}

func MakeBatchWithRootAndCommitments(
	batch *models.BatchWithAccountRoot,
	submissionBlock uint32,
	commitments []models.CommitmentWithTokenID,
) *BatchWithRootAndCommitments {
	return &BatchWithRootAndCommitments{
		Batch:           *MakeBatch(&batch.Batch, submissionBlock),
		AccountTreeRoot: batch.AccountTreeRoot,
		Commitments:     commitments,
	}
}
