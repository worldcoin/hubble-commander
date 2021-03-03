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
	leafPath := &models.MerklePath{
		Path:  index,
		Depth: 32,
	}
	rootPath := &models.MerklePath{
		Path:  0,
		Depth: 0,
	}
	prevLeaf, err := s.storage.GetStateNodeByPath(leafPath)
	if err != nil {
		return err
	}
	prevRoot, err := s.storage.GetStateNodeByPath(rootPath)
	if err != nil {
		return err
	}

	leaf, err := NewStateLeaf(state)
	if err != nil {
		return err
	}

	err = s.storage.AddStateLeaf(leaf)
	if err != nil {
		return err
	}

	witnessPaths, err := leafPath.GetWitnessPaths()
	if err != nil {
		return err
	}

	currentPath := leafPath
	currentHash := leaf.DataHash
	for _, witnessPath := range witnessPaths {
		err = s.storage.AddOrUpdateStateNode(&models.StateNode{
			MerklePath: *currentPath,
			DataHash:   currentHash,
		})
		if err != nil {
			return err
		}

		// nolint:gosec,govet
		witness, err := s.storage.GetStateNodeByPath(&witnessPath)
		if err != nil {
			return err
		}

		if currentPath.IsLeftNode() {
			currentHash = hashTwo(currentHash, witness.DataHash)
		} else {
			currentHash = hashTwo(witness.DataHash, currentHash)
		}
		currentPath, err = currentPath.Parent()
		if err != nil {
			return err
		}
	}

	err = s.storage.AddOrUpdateStateNode(&models.StateNode{
		MerklePath: *currentPath,
		DataHash:   currentHash,
	})
	if err != nil {
		return err
	}

	update := &models.StateUpdate{
		MerklePath:  *leafPath,
		CurrentHash: leaf.DataHash,
		CurrentRoot: currentHash,
		PrevHash:    prevLeaf.DataHash,
		PrevRoot:    prevRoot.DataHash,
	}
	err = s.storage.AddStateUpdate(update)
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
