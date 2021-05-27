package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetBatches(from, to *models.Uint256) ([]models.BatchWithSubmissionBlock, error) {
	batches, err := a.storage.GetBatchesInRange(from, to)
	if err != nil {
		return []models.BatchWithSubmissionBlock{}, err
	}

	blocksToFinalise, err := a.client.GetBlocksToFinalise()
	if err != nil {
		return nil, err
	}

	batchesWithSubmission := make([]models.BatchWithSubmissionBlock, 0, len(batches))
	for i := range batches {
		batchesWithSubmission = append(batchesWithSubmission, models.BatchWithSubmissionBlock{
			Batch:           batches[i],
			SubmissionBlock: *batches[i].FinalisationBlock - uint32(*blocksToFinalise),
		})
	}
	return batchesWithSubmission, nil
}
