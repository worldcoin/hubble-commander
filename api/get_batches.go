package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/storage"
)

var getBatchesAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(30001, "batches not found"),
}

func (a *API) GetBatches(from, to *models.Uint256) ([]interface{}, error) {
	batches, err := a.unsafeGetBatches(from, to)
	if err != nil {
		return nil, sanitizeError(err, getBatchesAPIErrors)
	}

	return batches, nil
}

func (a *API) unsafeGetBatches(from, to *models.Uint256) ([]interface{}, error) {
	batches, err := a.storage.GetBatchesInRange(from, to)
	if err != nil {
		return []interface{}{}, err
	}

	batchesWithSubmission := make([]interface{}, 0, len(batches))
	for i := range batches {
		status := calculateBatchStatus(a.storage.GetLatestBlockNumber(), &batches[i])

		if *status == batchstatus.Submitted {
			batchesWithSubmission = append(batchesWithSubmission, dto.MakeSubmittedBatch(&batches[i]))
		} else {
			submissionBlock, err := a.getSubmissionBlock(&batches[i])
			if err != nil {
				return nil, err
			}

			batchesWithSubmission = append(batchesWithSubmission, dto.MakeBatch(&batches[i], submissionBlock, status))
		}
	}
	return batchesWithSubmission, nil
}
