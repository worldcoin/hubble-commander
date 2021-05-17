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

const leafDepth = 32

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

func (s *StateTree) LeafNode(id uint32) (*models.StateNode, error) {
	leafPath := &models.MerklePath{
		Path:  id,
		Depth: leafDepth,
	}
	return s.storage.GetStateNodeByPath(leafPath)
}

func (s *StateTree) Leaf(id uint32) (*models.StateLeaf, error) {
	leafPath := &models.MerklePath{
		Path:  id,
		Depth: leafDepth,
	}
	return s.storage.GetStateLeafByPath(leafPath)
}

func (s *StateTree) Set(id uint32, state *models.UserState) (err error) {
	tx, storage, err := s.storage.BeginTransaction(TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return
	}
	defer tx.Rollback(&err)

	err = NewStateTree(storage).unsafeSet(id, state)
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

		currentRootHash, err = s.updateStateNodes(&latestStateUpdate.StateID, &latestStateUpdate.PrevHash)
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
		StateID:     prevLeaf.MerklePath,
		CurrentHash: currentLeaf.DataHash,
		CurrentRoot: *currentRoot,
		PrevHash:    prevLeaf.DataHash,
		PrevRoot:    *prevRoot,
	})
}

func nodesSliceToMap(nodes []models.StateNode) map[models.MerklePath]common.Hash {
	result := make(map[models.MerklePath]common.Hash, len(nodes))
	for i := range nodes {
		result[nodes[i].MerklePath] = nodes[i].DataHash
	}
	return result
}

func (s *StateTree) updateStateNodes(leafPath *models.MerklePath, newLeafHash *common.Hash) (*common.Hash, error) {
	witnessPaths, err := leafPath.GetWitnessPaths()
	if err != nil {
		return nil, err
	}

	nodes, err := s.storage.GetStateNodes(witnessPaths)
	if err != nil {
		return nil, err
	}

	nodesMap := nodesSliceToMap(nodes)
	nodes = make([]models.StateNode, 0, len(witnessPaths))
	currentHash := *newLeafHash
	var currentPath *models.MerklePath
	for _, witnessPath := range witnessPaths {
		currentPath, err = witnessPath.Sibling()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, models.StateNode{
			MerklePath: *currentPath,
			DataHash:   currentHash,
		})
		currentHash = s.calculateParentHash(&currentHash, currentPath, getWitnessHash(nodesMap, witnessPath))
	}

	nodes = append(nodes, models.StateNode{
		MerklePath: rootPath,
		DataHash:   currentHash,
	})

	err = s.storage.BatchUpsertStateNodes(nodes)
	if err != nil {
		return nil, err
	}

	return &currentHash, nil
}

func getWitnessHash(nodes map[models.MerklePath]common.Hash, path models.MerklePath) common.Hash {
	witnessHash, ok := nodes[path]
	if !ok {
		return GetZeroHash(leafDepth - uint(path.Depth))
	}
	return witnessHash
}

func (s *StateTree) calculateParentHash(
	currentHash *common.Hash,
	currentPath *models.MerklePath,
	witnessHash common.Hash,
) common.Hash {
	if currentPath.IsLeftNode() {
		return utils.HashTwo(*currentHash, witnessHash)
	}

	return utils.HashTwo(witnessHash, *currentHash)
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
		PubkeyID: big.NewInt(int64(state.PubKeyID)),
		TokenID:  &state.TokenIndex.Int,
		Balance:  &state.Balance.Int,
		Nonce:    &state.Nonce.Int,
	}
}
