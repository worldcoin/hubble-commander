package storage

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const StateTreeDepth = merkletree.MaxDepth

var stateUpdatePrefix = []byte("bh_" + reflect.TypeOf(models.StateUpdate{}).Name())

type StateTree struct {
	storageBase *StorageBase
	merkleTree  *StoredMerkleTree
}

func NewStateTree(storageBase *StorageBase) *StateTree {
	return &StateTree{
		storageBase: storageBase,
		merkleTree:  NewStoredMerkleTree("state", storageBase.Badger),
	}
}

func (s *StateTree) Root() (*common.Hash, error) {
	return s.merkleTree.Root()
}

func (s *StateTree) LeafNode(stateID uint32) (*models.MerkleTreeNode, error) {
	return s.merkleTree.Get(models.MerklePath{
		Path:  stateID,
		Depth: StateTreeDepth,
	})
}

func (s *StateTree) Leaf(stateID uint32) (*models.StateLeaf, error) {
	leaf, err := s.storageBase.GetStateLeaf(stateID)
	if IsNotFoundError(err) {
		return &models.StateLeaf{
			StateID:  stateID,
			DataHash: merkletree.GetZeroHash(0),
		}, nil
	} else if err != nil {
		return nil, err
	}
	return leaf, nil
}

// Set returns a witness containing 32 elements for the current set operation
func (s *StateTree) Set(id uint32, state *models.UserState) (models.Witness, error) {
	tx, storage, err := s.storageBase.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	witness, err := NewStateTree(storage).unsafeSet(id, state)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return witness, nil
}

func (s *StateTree) GetWitness(stateID uint32) (models.Witness, error) {
	return s.merkleTree.GetWitness(models.MakeMerklePathFromLeafID(stateID))
}

func (s *StateTree) RevertTo(targetRootHash common.Hash) error {
	txn, storage, err := s.storageBase.BeginTransaction(TxOptions{Badger: true})
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

			currentRootHash, err = stateTree.revertState(stateUpdate)
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

func decodeStateUpdate(item *bdg.Item) (*models.StateUpdate, error) {
	var stateUpdate models.StateUpdate
	err := item.Value(func(v []byte) error {
		return badger.Decode(v, &stateUpdate)
	})
	if err != nil {
		return nil, err
	}
	err = badger.DecodeKey(item.Key(), &stateUpdate.ID, stateUpdatePrefix)
	if err != nil {
		return nil, err
	}
	return &stateUpdate, nil
}

func (s *StateTree) unsafeSet(index uint32, state *models.UserState) (models.Witness, error) {
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

	err = s.storageBase.UpsertStateLeaf(currentLeaf)
	if err != nil {
		return nil, err
	}

	prevLeafPath := models.MakeMerklePathFromLeafID(prevLeaf.StateID)
	currentRoot, witness, err := s.merkleTree.SetNode(&prevLeafPath, currentLeaf.DataHash)
	if err != nil {
		return nil, err
	}

	err = s.storageBase.AddStateUpdate(&models.StateUpdate{
		CurrentRoot:   *currentRoot,
		PrevRoot:      *prevRoot,
		PrevStateLeaf: *prevLeaf,
	})
	if err != nil {
		return nil, err
	}

	return witness, nil
}

func (s *StateTree) revertState(stateUpdate *models.StateUpdate) (*common.Hash, error) {
	err := s.storageBase.UpsertStateLeaf(&stateUpdate.PrevStateLeaf)
	if err != nil {
		return nil, err
	}

	leafPath := models.MakeMerklePathFromLeafID(stateUpdate.PrevStateLeaf.StateID)
	currentRootHash, _, err := s.merkleTree.SetNode(&leafPath, stateUpdate.PrevStateLeaf.DataHash)
	if err != nil {
		return nil, err
	}
	if *currentRootHash != stateUpdate.PrevRoot {
		return nil, fmt.Errorf("unexpected state root after state update rollback")
	}

	err = s.storageBase.DeleteStateUpdate(stateUpdate.ID)
	if err != nil {
		return nil, err
	}

	return currentRootHash, nil
}

func (s *StateTree) getMerkleTreeNodeByPath(path *models.MerklePath) (*models.MerkleTreeNode, error) {
	node, err := s.merkleTree.Get(*path)
	if err != nil {
		return nil, err
	}
	return node, nil
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
		TokenID:  state.TokenID.ToBig(),
		Balance:  state.Balance.ToBig(),
		Nonce:    state.Nonce.ToBig(),
	}
}
