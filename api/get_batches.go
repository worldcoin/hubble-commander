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
		if batches[i].Hash == nil {
			continue
		}
		submissionBlock := *batches[i].FinalisationBlock - uint32(*blocksToFinalise)
		batchesWithSubmission = append(batchesWithSubmission, *dto.MakeBatch(&models.Batch{
			ID:                batches[i].ID,
			Type:              batches[i].Type,
			TransactionHash:   batches[i].TransactionHash,
			Hash:              batches[i].Hash,
			BlockTime:         batches[i].BlockTime,
			FinalisationBlock: batches[i].FinalisationBlock,
		}, submissionBlock))
	}
	return batchesWithSubmission, nil
}
