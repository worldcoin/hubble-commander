package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const (
	StateTreeDepth = merkletree.MaxDepth
	StateTreeSize  = int64(1) << StateTreeDepth
)

type StateTree struct {
	database   *Database
	merkleTree *StoredMerkleTree
}

func NewStateTree(database *Database) *StateTree {
	return &StateTree{
		database:   database,
		merkleTree: NewStoredMerkleTree("state", database, StateTreeDepth),
	}
}

func (s *StateTree) copyWithNewDatabase(database *Database) *StateTree {
	return NewStateTree(database)
}

func (s *StateTree) Root() (*common.Hash, error) {
	return s.merkleTree.Root()
}

func (s *StateTree) Leaf(stateID uint32) (stateLeaf *models.StateLeaf, err error) {
	var storedLeaf stored.StateLeaf
	err = s.database.Badger.Get(stateID, &storedLeaf)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("state leaf"))
	}
	if err != nil {
		return nil, err
	}
	return storedLeaf.ToModelsStateLeaf(), nil
}

func (s *StateTree) LeafOrEmpty(stateID uint32) (*models.StateLeaf, error) {
	leaf, err := s.Leaf(stateID)
	if IsNotFoundError(err) {
		return &models.StateLeaf{
			StateID:  stateID,
			DataHash: merkletree.GetZeroHash(0),
		}, nil
	}
	return leaf, err
}

func (s *StateTree) NextAvailableStateID() (*uint32, error) {
	return s.NextVacantSubtree(0)
}

// NextVacantSubtree returns the starting index of a vacant subtree of at least `subtreeDepth`.
// `subtreeDepth` can be set to 0 to only search for a single vacant node.
func (s *StateTree) NextVacantSubtree(subtreeDepth uint8) (*uint32, error) {
	subtreeWidth := int64(1) << subtreeDepth // Number of leaves in the subtree.

	prevTakenNodeIndex := int64(-1)
	result := uint32(0)

	// The iterator will scan over the state tree left-to-right detecting any gaps along the way.
	// If a gap is detected its checked if its suitable for the given subtree regarding both alignment and size.
	// An iterator will return the index of the first such gap it detects.
	err := s.database.Badger.Iterator(stored.StateLeafPrefix, db.KeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		var key uint32
		err := db.DecodeKey(item.Key(), &key, stored.StateLeafPrefix)
		if err != nil {
			return false, err
		}
		currentNodeIndex := int64(key)

		if currentNodeIndex != prevTakenNodeIndex+1 { // We detected a gap
			roundedNodeIndex := roundAndValidateStateTreeSlot(prevTakenNodeIndex+1, currentNodeIndex, subtreeWidth)
			if roundedNodeIndex != nil {
				result = uint32(*roundedNodeIndex)
				return true, nil
			}
		}

		prevTakenNodeIndex = currentNodeIndex
		return false, nil
	})
	if err == db.ErrIteratorFinished { // We finished without finding any gaps, try to append the subtree at the end.
		roundedNodeIndex := roundAndValidateStateTreeSlot(prevTakenNodeIndex+1, StateTreeSize, subtreeWidth)
		if roundedNodeIndex == nil {
			return nil, errors.WithStack(NewNoVacantSubtreeError(subtreeDepth))
		}
		return ref.Uint32(uint32(*roundedNodeIndex)), nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &result, nil
}

func roundAndValidateStateTreeSlot(rangeStart, rangeEnd, subtreeWidth int64) *int64 {
	// Check if we are aligned
	roundedNodeIndex := rangeStart
	if roundedNodeIndex%subtreeWidth != 0 {
		// If its not aligned to subtree size, round it to the next slot
		roundedNodeIndex += subtreeWidth - roundedNodeIndex%subtreeWidth
	}

	// Check if we fit in the current gap
	if roundedNodeIndex+subtreeWidth > rangeEnd {
		// Can't fit in the current gap
		return nil
	}

	return ref.Int64(roundedNodeIndex)
}

// Set returns a witness containing 32 elements for the current set operation
func (s *StateTree) Set(id uint32, state *models.UserState) (witness models.Witness, err error) {
	err = s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		witness, err = NewStateTree(txDatabase).unsafeSet(id, state)
		return err
	})
	if err != nil {
		return nil, err
	}

	return witness, nil
}

func (s *StateTree) GetLeafWitness(stateID uint32) (models.Witness, error) {
	return s.merkleTree.GetWitness(models.MakeMerklePathFromLeafID(stateID))
}

func (s *StateTree) GetNodeWitness(path models.MerklePath) (models.Witness, error) {
	return s.merkleTree.GetWitness(path)
}

func (s *StateTree) RevertTo(targetRootHash common.Hash) error {
	currentRootHash, err := s.Root()
	if err != nil {
		return err
	}
	if *currentRootHash == targetRootHash {
		return nil
	}

	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) (err error) {
		stateTree := NewStateTree(txDatabase)

		err = txDatabase.Badger.Iterator(models.StateUpdatePrefix, db.ReversePrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
			var stateUpdate *models.StateUpdate
			stateUpdate, err = decodeStateUpdate(item)
			if err != nil {
				return false, err
			}
			if stateUpdate.CurrentRoot != *currentRootHash {
				panic("invalid current root of a previous state update, this should never happen")
			}

			currentRootHash, err = stateTree.revertState(stateUpdate)
			if err != nil {
				return false, err
			}
			return *currentRootHash == targetRootHash, nil
		})
		if err != nil && err != db.ErrIteratorFinished {
			return errors.WithStack(err)
		}

		if *currentRootHash != targetRootHash {
			return errors.WithStack(ErrNonexistentState)
		}
		return nil
	})
}

func decodeStateUpdate(item *bdg.Item) (*models.StateUpdate, error) {
	var stateUpdate models.StateUpdate
	err := item.Value(func(v []byte) error {
		return db.Decode(v, &stateUpdate)
	})
	if err != nil {
		return nil, err
	}
	err = db.DecodeKey(item.Key(), &stateUpdate.ID, models.StateUpdatePrefix)
	if err != nil {
		return nil, err
	}
	return &stateUpdate, nil
}

func (s *StateTree) unsafeSet(index uint32, state *models.UserState) (models.Witness, error) {
	prevLeaf, err := s.LeafOrEmpty(index)
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

	err = s.upsertStateLeaf(currentLeaf)
	if err != nil {
		return nil, err
	}

	prevLeafPath := models.MakeMerklePathFromLeafID(prevLeaf.StateID)
	currentRoot, witness, err := s.merkleTree.SetNode(&prevLeafPath, currentLeaf.DataHash)
	if err != nil {
		return nil, err
	}

	err = s.addStateUpdate(&models.StateUpdate{
		CurrentRoot:   *currentRoot,
		PrevRoot:      *prevRoot,
		PrevStateLeaf: *prevLeaf,
	})
	if err != nil {
		return nil, err
	}

	return witness, nil
}

func (s *StateTree) getLeafByPubKeyIDAndTokenID(pubKeyID uint32, tokenID models.Uint256) (*models.StateLeaf, error) {
	stateLeaves := make([]stored.StateLeaf, 0, 1)
	err := s.database.Badger.Find(
		&stateLeaves,
		bh.Where("PubKeyID").Eq(pubKeyID).Index("PubKeyID"),
	)
	if err != nil {
		return nil, err
	}
	if len(stateLeaves) == 0 {
		return nil, errors.WithStack(NewNotFoundError("state leaf"))
	}

	for i := range stateLeaves {
		if stateLeaves[i].TokenID.Eq(&tokenID) {
			return stateLeaves[i].ToModelsStateLeaf(), nil
		}
	}

	return nil, errors.WithStack(NewNotFoundError("state leaf"))
}

func (s *StateTree) revertState(stateUpdate *models.StateUpdate) (*common.Hash, error) {
	err := s.upsertStateLeaf(&stateUpdate.PrevStateLeaf)
	if err != nil {
		return nil, err
	}

	leafPath := models.MakeMerklePathFromLeafID(stateUpdate.PrevStateLeaf.StateID)
	currentRootHash, _, err := s.merkleTree.SetNode(&leafPath, stateUpdate.PrevStateLeaf.DataHash)
	if err != nil {
		return nil, err
	}
	if *currentRootHash != stateUpdate.PrevRoot {
		panic("unexpected state root after state update rollback, this should never happen")
	}

	err = s.deleteStateUpdate(stateUpdate.ID)
	if err != nil {
		return nil, err
	}

	return currentRootHash, nil
}

func (s *StateTree) IterateLeaves(action func(stateLeaf *models.StateLeaf) error) error {
	err := s.database.Badger.Iterator(stored.StateLeafPrefix, db.PrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
		var stateLeaf stored.StateLeaf
		err := item.Value(stateLeaf.SetBytes)
		if err != nil {
			return false, err
		}

		err = action(stateLeaf.ToModelsStateLeaf())
		return false, err
	})
	if err != nil && err != db.ErrIteratorFinished {
		return err
	}
	return nil
}

func NewStateLeaf(stateID uint32, state *models.UserState) (*models.StateLeaf, error) {
	dataHash, err := encoder.HashUserState(state)
	if err != nil {
		return nil, err
	}
	return &models.StateLeaf{
		StateID:   stateID,
		DataHash:  *dataHash,
		UserState: *state,
	}, nil
}
