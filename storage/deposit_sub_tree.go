package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *DepositStorage) AddPendingDepositSubtree(subtree *models.PendingDepositSubtree) error {
	return s.database.Badger.Upsert(subtree.ID, *subtree)
}

func (s *DepositStorage) GetPendingDepositSubtree(subtreeID models.Uint256) (*models.PendingDepositSubtree, error) {
	var subtree models.PendingDepositSubtree
	err := s.database.Badger.Get(subtreeID, &subtree)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("deposit sub tree"))
	}
	if err != nil {
		return nil, err
	}

	subtree.ID = subtreeID

	return &subtree, nil
}

func (s *DepositStorage) GetFirstPendingDepositSubtree() (subtree *models.PendingDepositSubtree, err error) {
	err = s.database.Badger.Iterator(models.PendingDepositSubtreePrefix, db.KeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		subtree, err = decodePendingDepositSubtree(item)
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err == db.ErrIteratorFinished {
		return nil, errors.WithStack(NewNotFoundError("deposit sub tree"))
	}
	return subtree, err
}

func (s *DepositStorage) RemovePendingDepositSubtrees(subtreeIDs ...models.Uint256) error {
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		for i := range subtreeIDs {
			err := txDatabase.Badger.Delete(subtreeIDs[i], models.PendingDepositSubtree{})
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

func decodePendingDepositSubtree(item *bdg.Item) (*models.PendingDepositSubtree, error) {
	var subtree models.PendingDepositSubtree
	err := item.Value(subtree.SetBytes)
	if err != nil {
		return nil, err
	}

	err = db.DecodeKey(item.Key(), &subtree.ID, models.PendingDepositSubtreePrefix)
	if err != nil {
		return nil, err
	}
	return &subtree, nil
}
