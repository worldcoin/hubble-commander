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

// reads the pending balances from badger
func (a *API) GetPendingPubkeyBalances(ctx context.Context, startPrefix []byte, pageSize uint32) ([]dto.PubkeyBalance, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	return a.storage.GetPendingPubkeyBalances(startPrefix, pageSize)
}

// scans the mempool to recompute the pending balances
func (a *API) RecomputePubkeyBalances(ctx context.Context, startPrefix []byte, pageSize uint32) ([]dto.PubkeyBalance, error) {
	err := a.verifyAuthKey(ctx)
	if err != nil {
		return nil, err
	}

	return a.storage.RecomputePendingPubkeyBalances(startPrefix, pageSize)
}
