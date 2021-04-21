package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetBatches(from, to *models.Uint256) ([]models.Batch, error) {
	batches, err := a.storage.GetBatchesInRange(from, to)
	if err != nil {
		return nil, err
	}

	return batches, nil
}
