package storage

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

var ErrRanOutOfPendingDeposits = fmt.Errorf(
	"the commander ran out of the pending deposits for already built deposit sub trees on chain" +
		" - this should never happen",
)

type DepositStorage struct {
	database *Database
}

func NewDepositStorage(database *Database) *DepositStorage {
	return &DepositStorage{
		database: database,
	}
}

func (s *DepositStorage) copyWithNewDatabase(database *Database) *DepositStorage {
	newDepositStorage := *s
	newDepositStorage.database = database

	return &newDepositStorage
}

func (s *DepositStorage) AddPendingDeposit(deposit *models.PendingDeposit) error {
	return s.database.Badger.Upsert(deposit.ID, *deposit)
}

func (s *DepositStorage) RemovePendingDeposits(deposits []models.PendingDeposit) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range deposits {
		err := txDatabase.Badger.Delete(deposits[i].ID, models.PendingDeposit{})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DepositStorage) GetFirstPendingDeposits(amount int) ([]models.PendingDeposit, error) {
	deposits := make([]models.PendingDeposit, 0, amount)
	keyIteratorOpts := bdg.IteratorOptions{
		PrefetchSize: amount, // prefetch pending deposits for performance
	}
	err := s.database.Badger.Iterator(models.PendingDepositPrefix, keyIteratorOpts, func(item *bdg.Item) (bool, error) {
		deposit, err := decodeDeposit(item)
		if err != nil {
			return false, err
		}
		deposits = append(deposits, *deposit)
		return len(deposits) == amount, nil
	})
	if err == db.ErrIteratorFinished {
		return nil, errors.WithStack(ErrRanOutOfPendingDeposits)
	}
	if err != nil {
		return nil, err
	}
	return deposits, nil
}

func decodeDeposit(item *bdg.Item) (*models.PendingDeposit, error) {
	var deposit models.PendingDeposit
	err := item.Value(deposit.SetBytes)
	if err != nil {
		return nil, err
	}
	return &deposit, nil
}
