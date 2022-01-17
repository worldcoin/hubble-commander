package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

var getBatchAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(30000, "batch not found"),
}

func (a *API) GetBatchByHash(hash common.Hash) (*dto.BatchWithRootAndCommitments, error) {
	batch, err := a.unsafeGetBatchByHash(hash)
	if err != nil {
		return nil, sanitizeError(err, getBatchAPIErrors)
	}

	return batch, nil
}

func (a *API) unsafeGetBatchByHash(hash common.Hash) (*dto.BatchWithRootAndCommitments, error) {
	batch, err := a.storage.GetBatchByHash(hash)
	if err != nil {
		return nil, err
	}
	submissionBlock, err := a.getSubmissionBlock(batch)
	if err != nil {
		return nil, err
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(batch.ID)
	if err != nil {
		return nil, err
	}

	return a.createBatchWithCommitments(batch, submissionBlock, commitments)
}

func (a *API) GetBatchByID(id models.Uint256) (*dto.BatchWithRootAndCommitments, error) {
	batch, err := a.unsafeGetBatchByID(id)
	if err != nil {
		return nil, sanitizeError(err, getBatchAPIErrors)
	}

	return batch, nil
}

func (a *API) unsafeGetBatchByID(id models.Uint256) (*dto.BatchWithRootAndCommitments, error) {
	batch, err := a.storage.GetMinedBatch(id)
	if err != nil {
		return nil, err
	}
	submissionBlock, err := a.getSubmissionBlock(batch)
	if err != nil {
		return nil, err
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(batch.ID)
	if err != nil {
		return nil, err
	}

	return a.createBatchWithCommitments(batch, submissionBlock, commitments)
}

func (a *API) createBatchWithCommitments(
	batch *models.Batch,
	submissionBlock uint32,
	commitments []models.Commitment,
) (*dto.BatchWithRootAndCommitments, error) {
	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return a.createBatchWithTxCommitments(batch, submissionBlock, commitments)
	case batchtype.MassMigration:
		return a.createBatchWithMMCommitments(batch, submissionBlock, commitments)
	default:
		return &dto.BatchWithRootAndCommitments{
			Batch: dto.Batch{
				ID:                batch.ID,
				Hash:              batch.Hash,
				Type:              batch.Type,
				TransactionHash:   batch.TransactionHash,
				SubmissionBlock:   submissionBlock,
				SubmissionTime:    batch.SubmissionTime,
				FinalisationBlock: batch.FinalisationBlock,
			},
			AccountTreeRoot: batch.AccountTreeRoot,
			Commitments:     nil,
		}, nil
	}
}

func (a *API) createBatchWithTxCommitments(
	batch *models.Batch,
	submissionBlock uint32,
	commitments []models.Commitment,
) (*dto.BatchWithRootAndCommitments, error) {
	batchCommitments := make([]dto.CommitmentWithTokenID, 0, len(commitments))
	for i := range commitments {
		stateLeaf, err := a.storage.StateTree.Leaf(commitments[i].ToTxCommitment().FeeReceiver)
		if err != nil {
			return nil, err
		}

		batchCommitments = append(batchCommitments, dto.MakeCommitmentWithTokenIDFromTxCommitment(
			commitments[i].ToTxCommitment(),
			stateLeaf.TokenID,
		))
	}
	return dto.MakeBatchWithRootAndCommitments(batch, submissionBlock, batchCommitments), nil
}

func (a *API) createBatchWithMMCommitments(
	batch *models.Batch,
	submissionBlock uint32,
	commitments []models.Commitment,
) (*dto.BatchWithRootAndCommitments, error) {
	batchCommitments := make([]dto.CommitmentWithTokenID, 0, len(commitments))
	for i := range commitments {
		stateLeaf, err := a.storage.StateTree.Leaf(commitments[i].ToMMCommitment().FeeReceiver)
		if err != nil {
			return nil, err
		}

		batchCommitments = append(batchCommitments, dto.MakeCommitmentWithTokenIDFromMMCommitment(
			commitments[i].ToMMCommitment(),
			stateLeaf.TokenID,
		))
	}
	return dto.MakeBatchWithRootAndCommitments(batch, submissionBlock, batchCommitments), nil
}

func (a *API) getSubmissionBlock(batch *models.Batch) (uint32, error) {
	if batch.ID.IsZero() {
		return *batch.FinalisationBlock, nil
	}

	blocks, err := a.client.GetBlocksToFinalise()
	if err != nil {
		return 0, err
	}
	return *batch.FinalisationBlock - uint32(*blocks), nil
}
