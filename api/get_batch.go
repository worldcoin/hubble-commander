package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

var getBatchAPIErrors = map[error]*APIError{
	&storage.NotFoundError{}: NewAPIError(30000, "batch not found"),
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
	submissionBlock, err := a.getSubmissionBlock(*batch.FinalisationBlock)
	if err != nil {
		return nil, err
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(batch.ID)
	if err != nil {
		return nil, err
	}
	return createBatchWithCommitments(batch, submissionBlock, commitments)
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
	submissionBlock, err := a.getSubmissionBlock(*batch.FinalisationBlock)
	if err != nil {
		return nil, err
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(batch.ID)
	if err != nil {
		return nil, err
	}
	return createBatchWithCommitments(batch, submissionBlock, commitments)
}

func createBatchWithCommitments(
	batch *models.Batch,
	submissionBlock uint32,
	commitments []models.CommitmentWithTokenID,
) (*dto.BatchWithRootAndCommitments, error) {
	for i := range commitments {
		commitments[i].LeafHash = commitments[i].CalcLeafHash(batch.AccountTreeRoot)
	}
	return dto.MakeBatchWithRootAndCommitments(batch, submissionBlock, commitments), nil
}

func (a *API) getSubmissionBlock(finalisationBlock uint32) (uint32, error) {
	blocks, err := a.client.GetBlocksToFinalise()
	if err != nil {
		return 0, err
	}
	return finalisationBlock - uint32(*blocks), nil
}
