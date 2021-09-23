package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
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
