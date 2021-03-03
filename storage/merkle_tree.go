package storage

import (
	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	rootPath = models.MerklePath{Path: 0, Depth: 0}
)

type StateTree struct {
	storage *Storage
}

func NewStateTree(storage *Storage) *StateTree {
	return &StateTree{storage}
}

func (s *StateTree) Root() (*common.Hash, error) {
	root, err := s.storage.GetStateNodeByPath(&rootPath)
	if err != nil {
		return nil, err
	}
	return &root.DataHash, nil
}

func (s *StateTree) LeafNode(index uint32) (*models.StateNode, error) {
	leafPath := &models.MerklePath{
		Path:  index,
		Depth: 32,
	}
	leaf, err := s.storage.GetStateNodeByPath(leafPath)
	if err != nil {
		return nil, err
	}
	return leaf, nil
}

func (s *StateTree) Set(index uint32, state *models.UserState) error {
	prevLeaf, err := s.LeafNode(index)
	if err != nil {
		return err
	}
	prevRoot, err := s.Root()
	if err != nil {
		return err
	}

	currentLeaf, err := NewStateLeaf(state)
	if err != nil {
		return err
	}

	err = s.storage.AddStateLeaf(currentLeaf)
	if err != nil {
		return err
	}

	currentRoot, err := s.updateStateNodes(&prevLeaf.MerklePath, &currentLeaf.DataHash)
	if err != nil {
		return err
	}

	err = s.storage.AddStateUpdate(&models.StateUpdate{
		MerklePath:  prevLeaf.MerklePath,
		CurrentHash: currentLeaf.DataHash,
		CurrentRoot: *currentRoot,
		PrevHash:    prevLeaf.DataHash,
		PrevRoot:    *prevRoot,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *StateTree) updateStateNodes(leafPath *models.MerklePath, newLeafHash *common.Hash) (*common.Hash, error) {
	witnessPaths, err := leafPath.GetWitnessPaths()
	if err != nil {
		return nil, err
	}

	currentHash := *newLeafHash
	for _, witnessPath := range witnessPaths {
		// nolint:govet
		currentPath, err := witnessPath.Sibling()
		if err != nil {
			return nil, err
		}

		err = s.storage.AddOrUpdateStateNode(&models.StateNode{
			MerklePath: *currentPath,
			DataHash:   currentHash,
		})
		if err != nil {
			return nil, err
		}

		currentHash, err = s.calculateParentHash(&currentHash, currentPath, witnessPath)
		if err != nil {
			return nil, err
		}
	}

	err = s.storage.AddOrUpdateStateNode(&models.StateNode{
		MerklePath: rootPath,
		DataHash:   currentHash,
	})
	if err != nil {
		return nil, err
	}

	return &currentHash, nil
}

func (s *StateTree) calculateParentHash(
	currentHash *common.Hash,
	currentPath *models.MerklePath,
	witnessPath models.MerklePath,
) (common.Hash, error) {
	witness, err := s.storage.GetStateNodeByPath(&witnessPath)
	if err != nil {
		return common.Hash{}, err
	}

	if currentPath.IsLeftNode() {
		return hashTwo(*currentHash, witness.DataHash), nil
	}

	return hashTwo(witness.DataHash, *currentHash), nil
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
