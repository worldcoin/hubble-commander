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

func (a *API) GetPendingStates(ctx context.Context, startStateID, pageSize uint32) ([]dto.UserStateWithID, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	return a.storage.GetPendingStates(startStateID, pageSize)
}

func (a *API) MempoolDropTransaction(ctx context.Context, stateID, nonce uint32) error {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return err
	}

	return a.storage.MempoolDropTransaction(stateID, nonce)
}
