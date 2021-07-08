package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
)

type StateTree2 struct {
	StateTree
	impl merkletree.StoredMerkleTree
}

func (s *StateTree2) Set(index uint32, state *models.UserState) (models.Witness, error) {
	// TODO start DB transaction

	prevLeaf, err := s.Leaf(index)
	if err != nil {
		return nil, err
	}

	prevRoot, err := s.Root()
	if err != nil {
		return nil, err
	}

	currentLeaf, err := NewStateLeaf(index, state)
	if err != nil {
		return nil, err
	}

	currentRoot, witness, err := s.impl.Set(index, currentLeaf)
	if err != nil {
		return nil, err
	}

	err = s.storage.AddStateUpdate(&models.StateUpdate{
		CurrentRoot:   *currentRoot,
		PrevRoot:      *prevRoot,
		PrevStateLeaf: *prevLeaf,
	})
	if err != nil {
		return nil, err
	}

	return witness, nil
}
