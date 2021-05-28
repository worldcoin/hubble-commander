package dto

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Batch struct {
	ID                *models.Uint256
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

func MakeBatch(batch *models.BatchWithSubmissionBlock) *Batch {
	return &Batch{
		ID:                batch.Number,
		Hash:              batch.Hash,
		Type:              batch.Type,
		TransactionHash:   batch.TransactionHash,
		SubmissionBlock:   batch.SubmissionBlock,
		FinalisationBlock: batch.FinalisationBlock,
	}
}

func MakeBatchWithRootAndCommitments(batch *models.BatchWithAccountRoot, commitments []models.CommitmentWithTokenID) *BatchWithRootAndCommitments {
	return &BatchWithRootAndCommitments{
		Batch:           *MakeBatch(&batch.BatchWithSubmissionBlock),
		AccountTreeRoot: batch.AccountTreeRoot,
		Commitments:     commitments,
	}
}
