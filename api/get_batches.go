package api

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *API) GetBatches(from, to *models.Uint256) ([]models.Batch, error) {
	return a.storage.GetBatchesInRange(from, to)
}
