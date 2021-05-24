package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"reflect"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const leafDepth = 32

var (
	rootPath          = models.MerklePath{Path: 0, Depth: 0}
	stateUpdatePrefix = []byte("bh_" + reflect.TypeOf(models.StateUpdate{}).Name())
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

func (s *StateTree) LeafNode(stateID uint32) (*models.StateNode, error) {
	leafPath := &models.MerklePath{
		Path:  stateID,
		Depth: leafDepth,
	}
	return s.storage.GetStateNodeByPath(leafPath)
}

func (s *StateTree) Leaf(stateID uint32) (*models.StateLeaf, error) {
	leaf, err := s.storage.GetStateLeaf(stateID)
	if IsNotFoundError(err) {
		return &models.StateLeaf{
			StateID:  stateID,
			DataHash: GetZeroHash(0),
		}, nil
	} else if err != nil {
		return nil, err
	}
	return leaf, nil
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
	txn, storage, err := s.storage.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer txn.Rollback(&err)

	stateTree := NewStateTree(storage)
	var currentRootHash *common.Hash
	err = storage.Badger.View(func(txn *bdg.Txn) error {
		currentRootHash, err = stateTree.Root()
		if err != nil {
			return err
		}

		opts := bdg.DefaultIteratorOptions
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()

		seekPrefix := make([]byte, 0, len(stateUpdatePrefix)+1)
		seekPrefix = append(seekPrefix, 0xFF)
		for it.Seek(seekPrefix); it.ValidForPrefix(stateUpdatePrefix); it.Next() {
			if *currentRootHash == targetRootHash {
				return nil
			}
			var stateUpdate *models.StateUpdate
			stateUpdate, err = decodeStateUpdate(it.Item())
			if err != nil {
				return err
			}
			if stateUpdate.CurrentRoot != *currentRootHash {
				panic("invalid current root of a previous state update, this should never happen")
			}

			err = storage.UpsertStateLeaf(&stateUpdate.PrevStateLeaf)
			if err != nil {
				return err
			}

			leafPath := models.MakeMerklePathFromStateID(stateUpdate.PrevStateLeaf.StateID)
			currentRootHash, err = stateTree.updateStateNodes(&leafPath, &stateUpdate.PrevStateLeaf.DataHash)
			if err != nil {
				return err
			}
			if *currentRootHash != stateUpdate.PrevRoot {
				return fmt.Errorf("unexpected state root after state update rollback")
			}

			err = storage.DeleteStateUpdate(stateUpdate.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if *currentRootHash != targetRootHash {
		return ErrNotExistentState
	}
	return txn.Commit()
}

func decodeKey(data []byte, key interface{}, prefix []byte) error {
	return gob.NewDecoder(bytes.NewReader(data[len(prefix):])).
		Decode(key)
}

func decodeStateUpdate(item *bdg.Item) (*models.StateUpdate, error) {
	var stateUpdate models.StateUpdate
	err := item.Value(func(v []byte) error {
		return badger.Decode(v, &stateUpdate)
	})
	if err != nil {
		return nil, err
	}
	err = decodeKey(item.Key(), &stateUpdate.ID, stateUpdatePrefix)
	if err != nil {
		return nil, err
	}
	return &stateUpdate, nil
}

func (s *StateTree) unsafeSet(index uint32, state *models.UserState) (err error) {
	prevLeaf, err := s.Leaf(index)
	if err != nil {
		return err
	}

	prevRoot, err := s.Root()
	if err != nil {
		return
	}

	currentLeaf, err := NewStateLeaf(index, state)
	if err != nil {
		return
	}

	err = s.storage.UpsertStateLeaf(currentLeaf)
	if err != nil {
		return
	}

	prevLeafPath := models.MakeMerklePathFromStateID(prevLeaf.StateID)
	currentRoot, err := s.updateStateNodes(&prevLeafPath, &currentLeaf.DataHash)
	if err != nil {
		return
	}

	return s.storage.AddStateUpdate(&models.StateUpdate{
		CurrentRoot:   *currentRoot,
		PrevRoot:      *prevRoot,
		PrevStateLeaf: *prevLeaf,
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

func NewStateLeaf(stateID uint32, state *models.UserState) (*models.StateLeaf, error) {
	encodedState, err := encoder.EncodeUserState(toContractUserState(state))
	if err != nil {
		return nil, err
	}
	dataHash := crypto.Keccak256Hash(encodedState)
	return &models.StateLeaf{
		StateID:   stateID,
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
