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
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		for i := range deposits {
			err := txDatabase.Badger.Delete(deposits[i].ID, models.PendingDeposit{})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DepositStorage) GetFirstPendingDeposits(amount int) ([]models.PendingDeposit, error) {
	deposits := make([]models.PendingDeposit, 0, amount)

	// set iterator opts to prefetch pending deposits for performance reasons
	keyIteratorOpts := bdg.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   amount,
	}

	err := s.database.Badger.Iterator(models.PendingDepositPrefix, keyIteratorOpts, func(item *bdg.Item) (bool, error) {
		deposit, err := decodeDeposit(item)
		if err != nil {
			return false, err
		}
		deposits = append(deposits, *deposit)
		return len(deposits) == amount, nil
	})
	if errors.Is(err, db.ErrIteratorFinished) {
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
