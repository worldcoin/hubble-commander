package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
)

func (s *StateTree) RevertToForDispute(targetRootHash common.Hash, invalidTransfer models.GenericTransfer) ([]models.StateMerkleProof, error) {
	txn, storage, err := s.storage.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return nil, err
	}
	defer txn.Rollback(&err)

	stateTree := NewStateTree(storage)
	var currentRootHash *common.Hash
	proofs := make([]models.StateMerkleProof, 0)

	userProofs, err := stateTree.createUserProofFromStateIDs(*invalidTransfer.GetToStateID(), invalidTransfer.GetFromStateID())
	if err != nil {
		return nil, err
	}
	proofs = append(proofs, userProofs...)

	err = storage.Badger.View(func(txn *bdg.Txn) error {
		currentRootHash, err = stateTree.Root()
		if err != nil {
			return err
		}

		opts := bdg.DefaultIteratorOptions
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()

		counter := 0
		prevStateLeaf := models.StateLeaf{}

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

			if counter%2 == 0 {
				prevStateLeaf = stateUpdate.PrevStateLeaf
			} else {
				userProofs, err = stateTree.createUserProofs(prevStateLeaf, stateUpdate.PrevStateLeaf)
				if err != nil {
					return err
				}
				proofs = append(proofs, userProofs...)
			}

			counter++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if *currentRootHash != targetRootHash {
		return nil, ErrNotExistentState
	}
	return proofs, txn.Commit()
}

func (s *StateTree) createUserProof(leaf *models.StateLeaf) (*models.StateMerkleProof, error) {
	witness, err := s.GetWitness(models.MakeMerklePathFromStateID(leaf.StateID))
	if err != nil {
		return nil, err
	}

	return &models.StateMerkleProof{
		UserState: &leaf.UserState,
		Witness:   witness,
	}, nil
}

func (s *StateTree) createUserProofFromStateIDs(stateIDs ...uint32) ([]models.StateMerkleProof, error) {
	proofs := make([]models.StateMerkleProof, 0, len(stateIDs))
	for i := range stateIDs {
		witness, err := s.GetWitness(models.MakeMerklePathFromStateID(stateIDs[i]))
		if err != nil {
			return nil, err
		}
		leaf, err := s.Leaf(stateIDs[i])
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, models.StateMerkleProof{
			UserState: &leaf.UserState,
			Witness:   witness,
		})
	}
	return proofs, nil
}

func (s *StateTree) createUserProofs(leaves ...models.StateLeaf) ([]models.StateMerkleProof, error) {
	proofs := make([]models.StateMerkleProof, 0, len(leaves))
	for i := range leaves {
		proof, err := s.createUserProof(&leaves[i])
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, *proof)
	}
	return proofs, nil
}

func (s *StateTree) GetWitness(leafPath models.MerklePath) (models.Witness, error) {
	witnessPaths, err := leafPath.GetWitnessPaths()
	if err != nil {
		return nil, err
	}
	nodes, err := s.storage.GetStateNodes(witnessPaths)
	if err != nil {
		return nil, err
	}
	witnesses := make([]common.Hash, 0, len(nodes))
	for i := range nodes {
		witnesses = append(witnesses, nodes[i].DataHash)
	}
	return witnesses, nil
}
