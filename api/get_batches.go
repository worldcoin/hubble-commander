package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetBatches(from, to *models.Uint256) ([]dto.Batch, error) {
	batches, err := a.storage.GetBatchesInRange(from, to)
	if err != nil {
		return []dto.Batch{}, err
	}

	blocksToFinalise, err := a.client.GetBlocksToFinalise()
	if err != nil {
		return nil, err
	}

	batchesWithSubmission := make([]dto.Batch, 0, len(batches))
	for i := range batches {
		submissionBlock := *batches[i].FinalisationBlock - uint32(*blocksToFinalise)
		batchesWithSubmission = append(batchesWithSubmission, *dto.MakeBatch(&models.BatchWithSubmissionBlock{
			Batch:           batches[i],
			SubmissionBlock: submissionBlock,
		}))
	}
	return batchesWithSubmission, nil
}
