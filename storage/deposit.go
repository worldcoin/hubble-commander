package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
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

// TODO-D is this needed?
func (s *DepositStorage) GetPendingDeposit(depositID *models.DepositID) (*models.PendingDeposit, error) {
	var deposit models.PendingDeposit
	err := s.database.Badger.Get(*depositID, &deposit)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("pending deposit"))
	}
	if err != nil {
		return nil, err
	}

	return &deposit, nil
}

func (s *DepositStorage) GetFirstPendingDeposits(amount int) ([]models.PendingDeposit, error) {
	deposits := make([]models.PendingDeposit, 0, amount)
	err := s.database.Badger.Iterator(models.PendingDepositPrefix, db.KeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		deposit, err := decodeDeposit(item)
		if err != nil {
			return false, err
		}
		deposits = append(deposits, *deposit)
		return len(deposits) == amount, nil
	})
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
