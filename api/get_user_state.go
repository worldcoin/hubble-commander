package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
)

func (a *API) GetUserState(id uint32) (*dto.UserStateWithID, error) {
	leaf, err := a.storage.StateTree.Leaf(id)
	if err != nil {
		return nil, err
	}
	return dto.NewUserStateWithID(leaf), nil
}
