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
	bh "github.com/timshannon/badgerhold/v4"
)

const StateTreeDepth = merkletree.MaxDepth

var stateUpdatePrefix = []byte("bh_" + reflect.TypeOf(models.StateUpdate{}).Name() + ":")

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
	var leaf models.FlatStateLeaf
	err = s.database.Badger.Get(stateID, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("state leaf")
	}
	if err != nil {
		return nil, err
	}
	return leaf.StateLeaf(), nil
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
	nextAvailableStateID := uint32(0)

	err := s.database.Badger.View(func(txn *bdg.Txn) error {
		opts := bdg.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()

		seekPrefix := make([]byte, 0, len(models.FlatStateLeafPrefix)+1)
		seekPrefix = append(seekPrefix, models.FlatStateLeafPrefix...)
		seekPrefix = append(seekPrefix, 0xFF) // Required to loop backwards

		it.Seek(seekPrefix)
		if it.ValidForPrefix(models.FlatStateLeafPrefix) {
			var key uint32
			err := badger.DecodeKey(it.Item().Key(), &key, models.FlatStateLeafPrefix)
			if err != nil {
				return err
			}
			nextAvailableStateID = key + 1
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &nextAvailableStateID, nil
}

// Set returns a witness containing 32 elements for the current set operation
func (s *StateTree) Set(id uint32, state *models.UserState) (models.Witness, error) {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	witness, err := NewStateTree(txDatabase).unsafeSet(id, state)
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
	txn, txDatabase, err := s.database.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer txn.Rollback(&err)

	stateTree := NewStateTree(txDatabase)
	var currentRootHash *common.Hash
	err = txDatabase.Badger.View(func(txn *bdg.Txn) error {
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
	leaves := make([]models.FlatStateLeaf, 0, 1)
	err := s.database.Badger.Find(
		&leaves,
		bh.Where("TokenID").Eq(tokenID).
			And("PubKeyID").Eq(pubKeyID).Index("PubKeyID"),
	)
	if err != nil {
		return nil, err
	}
	if len(leaves) == 0 {
		return nil, NewNotFoundError("state leaf")
	}
	return leaves[0].StateLeaf(), nil
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
		return nil, fmt.Errorf("unexpected state root after state update rollback")
	}

	err = s.deleteStateUpdate(stateUpdate.ID)
	if err != nil {
		return nil, err
	}

	return currentRootHash, nil
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
