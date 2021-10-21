package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
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

func (s *DepositStorage) DeletePendingDepositSubTrees(subTreeIDs ...models.Uint256) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range subTreeIDs {
		err = txDatabase.Badger.Delete(subTreeIDs[i], models.PendingDepositSubTree{})
		if err == bh.ErrNotFound {
			return errors.WithStack(NewNotFoundError("deposit sub tree"))
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
