package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type PendingStakeWithdrawalStorage struct {
	database *Database
}

func NewPendingStakeWithdrawalStorage(database *Database) *PendingStakeWithdrawalStorage {
	return &PendingStakeWithdrawalStorage{
		database: database,
	}
}

func (s *PendingStakeWithdrawalStorage) copyWithNewDatabase(database *Database) *PendingStakeWithdrawalStorage {
	newPendingStakeWithdrawalStorage := *s
	newPendingStakeWithdrawalStorage.database = database

	return &newPendingStakeWithdrawalStorage
}

func (s *PendingStakeWithdrawalStorage) AddPendingStakeWithdrawal(stake *models.PendingStakeWithdrawal) error {
	return s.database.Badger.Insert(stake.BatchID, stake)
}

func (s *PendingStakeWithdrawalStorage) RemovePendingStakeWithdrawal(batchID models.Uint256) error {
	var stake models.PendingStakeWithdrawal
	err := s.database.Badger.Delete(batchID, &stake)
	if err == bh.ErrNotFound {
		return errors.WithStack(NewNotFoundError("pending stake withdrawal"))
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *PendingStakeWithdrawalStorage) GetPendingStakeWithdrawalByFinalisationBlock(block uint32) (*models.PendingStakeWithdrawal, error) {
	var stake models.PendingStakeWithdrawal
	err := s.database.Badger.FindOne(&stake, bh.Where("FinalisationBlock").Eq(block))
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("pending stake withdrawal"))
	}
	return &stake, err
}
