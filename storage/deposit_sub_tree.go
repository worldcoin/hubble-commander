package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *DepositStorage) AddPendingDepositSubTree(subTree *models.PendingDepositSubTree) error {
	return s.database.Badger.Upsert(subTree.ID, *subTree)
}

// TODO - replace with a proper getter when implementing submit deposit batches
func (s *DepositStorage) GetPendingDepositSubTree(subTreeID models.Uint256) (*models.PendingDepositSubTree, error) {
	var subTree models.PendingDepositSubTree
	err := s.database.Badger.Get(subTreeID, &subTree)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("deposit sub tree"))
	}
	if err != nil {
		return nil, err
	}

	subTree.ID = subTreeID

	return &subTree, nil
}

func (s *DepositStorage) GetFirstPendingDepositSubTree() (subTree *models.PendingDepositSubTree, err error) {
	err = s.database.Badger.Iterator(models.PendingDepositSubTreePrefix, db.KeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		subTree, err = decodePendingDepositSubTree(item)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err == db.ErrIteratorFinished {
		return nil, errors.WithStack(NewNotFoundError("deposit sub tree"))
	}
	if err != nil {
		return nil, err
	}
	return subTree, nil
}

func (s *DepositStorage) DeletePendingDepositSubTrees(subTreeIDs ...models.Uint256) error {
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		for i := range subTreeIDs {
			err := txDatabase.Badger.Delete(subTreeIDs[i], models.PendingDepositSubTree{})
			if err == bh.ErrNotFound {
				return errors.WithStack(NewNotFoundError("deposit sub tree"))
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func decodePendingDepositSubTree(item *bdg.Item) (*models.PendingDepositSubTree, error) {
	var subTree models.PendingDepositSubTree
	err := item.Value(subTree.SetBytes)
	if err != nil {
		return nil, err
	}

	err = db.DecodeKey(item.Key(), &subTree.ID, models.PendingDepositSubTreePrefix)
	if err != nil {
		return nil, err
	}
	return &subTree, nil
}
