package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) RecomputePendingState(ctx context.Context, stateID uint32, mutate bool) (*dto.RecomputePendingState, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	return a.storage.RecomputePendingState(stateID, mutate)
}

func (a *API) GetPendingStates(ctx context.Context, startStateID uint32, pageSize uint32) ([]dto.UserStateWithID, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	return a.storage.GetPendingStates(startStateID, pageSize)
}
