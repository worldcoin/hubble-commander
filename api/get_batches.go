package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getBatchesAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(30001, "batches not found"),
}

func (a *API) GetBatches(from, to *models.Uint256) ([]dto.Batch, error) {
	batches, err := a.unsafeGetBatches(from, to)
	if err != nil {
		return nil, sanitizeError(err, getBatchesAPIErrors)
	}

	return batches, nil
}

func (a *API) unsafeGetBatches(from, to *models.Uint256) ([]dto.Batch, error) {
	batches, err := a.storage.GetBatchesInRange(from, to)
	if err != nil {
		return []dto.Batch{}, err
	}

	batchesWithSubmission := make([]dto.Batch, 0, len(batches))
	for i := range batches {
		if batches[i].Hash == nil {
			continue
		}

		submissionBlock, err := a.getSubmissionBlock(&batches[i])
		if err != nil {
			return nil, err
		}

		status := calculateBatchStatus(a.storage.GetLatestBlockNumber(), *batches[i].FinalisationBlock)

		batchesWithSubmission = append(batchesWithSubmission, *dto.MakeBatch(&batches[i], submissionBlock, status))
	}
	return batchesWithSubmission, nil
}
