package storage

import (
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
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
	return s.storage.GetStateNodeByPath(leafPath)
}

func (s *StateTree) Leaf(index uint32) (*models.StateLeaf, error) {
	leafPath := &models.MerklePath{
		Path:  index,
		Depth: 32,
	}
	return s.storage.GetStateLeafByPath(leafPath)
}

func (s *StateTree) Set(index uint32, state *models.UserState) (err error) {
	tx, storage, err := s.storage.BeginTransaction()
	if err != nil {
		return
	}
	defer tx.Rollback(&err)

	err = NewStateTree(storage).unsafeSet(index, state)
	if err != nil {
		return
	}

	return tx.Commit()
}

func (s *StateTree) RevertTo(targetRootHash common.Hash) error {
	_, err := s.storage.GetStateUpdateByRootHash(targetRootHash)
	if err != nil {
		if IsNotFoundError(err) {
			return fmt.Errorf("cannot revert to not existent state")
		}

		return err
	}

	stateTree := NewStateTree(s.storage)

	currentRootHash, err := stateTree.Root()
	if err != nil {
		return err
	}

	for *currentRootHash != targetRootHash {
		latestStateUpdate, err := s.storage.GetStateUpdateByRootHash(*currentRootHash)
		if err != nil {
			return err
		}

		currentRootHash, err = s.updateStateNodes(&latestStateUpdate.MerklePath, &latestStateUpdate.PrevHash)
		if err != nil {
			return err
		}
		if *currentRootHash != latestStateUpdate.PrevRoot {
			return fmt.Errorf("unexpected state root after state update rollback")
		}

		err = s.storage.DeleteStateUpdate(latestStateUpdate.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StateTree) unsafeSet(index uint32, state *models.UserState) (err error) {
	prevLeaf, err := s.LeafNode(index)
	if err != nil {
		return
	}
	prevRoot, err := s.Root()
	if err != nil {
		return
	}

	currentLeaf, err := NewStateLeaf(state)
	if err != nil {
		return
	}

	err = s.storage.AddStateLeaf(currentLeaf)
	if err != nil {
		return
	}

	currentRoot, err := s.updateStateNodes(&prevLeaf.MerklePath, &currentLeaf.DataHash)
	if err != nil {
		return
	}

	return s.storage.AddStateUpdate(&models.StateUpdate{
		MerklePath:  prevLeaf.MerklePath,
		CurrentHash: currentLeaf.DataHash,
		CurrentRoot: *currentRoot,
		PrevHash:    prevLeaf.DataHash,
		PrevRoot:    *prevRoot,
	})
}

func (s *StateTree) updateStateNodes(leafPath *models.MerklePath, newLeafHash *common.Hash) (*common.Hash, error) {
	witnessPaths, err := leafPath.GetWitnessPaths()
	if err != nil {
		return nil, err
	}

	currentHash := *newLeafHash
	for _, witnessPath := range witnessPaths {
		var currentPath *models.MerklePath
		currentPath, err = witnessPath.Sibling()
		if err != nil {
			return nil, err
		}

		err = s.storage.UpsertStateNode(&models.StateNode{
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

	err = s.storage.UpsertStateNode(&models.StateNode{
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
		return utils.HashTwo(*currentHash, witness.DataHash), nil
	}

	return utils.HashTwo(witness.DataHash, *currentHash), nil
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
		PubkeyID: big.NewInt(int64(state.PubkeyID)),
		TokenID:  &state.TokenIndex.Int,
		Balance:  &state.Balance.Int,
		Nonce:    &state.Nonce.Int,
	}
}
