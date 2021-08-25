package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
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

func (s *DepositStorage) AddDeposit(deposit *models.Deposit) error {
	return s.database.Badger.Upsert(deposit.ID, *deposit)
}

func (s *DepositStorage) GetDeposit(depositID *models.DepositID) (*models.Deposit, error) {
	var deposit models.Deposit
	err := s.database.Badger.Get(depositID, &deposit)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("deposit")
	}
	if err != nil {
		return nil, err
	}

	return &deposit, nil
}
