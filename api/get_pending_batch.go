package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetPendingBatches() ([]dto.Batch, error) {
	batches, err := a.unsafeGetPendingBatches()
	if err != nil {
		return nil, sanitizeError(err, map[error]*APIError{})
	}

	return batches, nil
}

func (a *API) unsafeGetPendingBatches() ([]dto.Batch, error) {
	batches, err := a.storage.GetPendingBatches()
	if err != nil {
		return nil, err
	}

	dtoBatches := make([]dto.Batch, 0, len(batches))
	for i := range batches {
		dtoBatches = append(dtoBatches, *dto.MakeBatch(&batches[i], 0))
	}
	return dtoBatches, nil
}
