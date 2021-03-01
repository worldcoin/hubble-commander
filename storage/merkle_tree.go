package storage

import (
	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/crypto"
)

type StateTree struct {
	storage *Storage
}

func NewStateTree(storage *Storage) *StateTree {
	return &StateTree{storage}
}

func (s *StateTree) Set(index uint32, state *models.UserState) error {
	leaf, err := NewStateLeaf(state)
	if err != nil {
		return err
	}

	err = s.storage.AddStateLeaf(leaf)
	if err != nil {
		return err
	}

	err = s.storage.AddStateNode(&models.StateNode{
		MerklePath: models.MerklePath{
			Path:  index,
			Depth: 32,
		},
		DataHash: leaf.DataHash,
	})
	if err != nil {
		return err
	}

	return nil
}

func NewStateLeaf(state *models.UserState) (*models.StateLeaf, error) {
	encodedState, err := encoder.EncodeUserState(toContractUserState(state))
	if err != nil {
		return nil, err
	}
	dataHash := crypto.Keccak256Hash(encodedState)
	return &models.StateLeaf{
		DataHash:  dataHash,
		UserState: *state,
	}, nil
}

func toContractUserState(state *models.UserState) generic.TypesUserState {
	return generic.TypesUserState{
		PubkeyID: &state.AccountIndex.Int,
		TokenID:  &state.TokenIndex.Int,
		Balance:  &state.Balance.Int,
		Nonce:    &state.Nonce.Int,
	}
}
