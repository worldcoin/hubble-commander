package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetPendingBatches(ctx context.Context) ([]dto.Batch, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

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
