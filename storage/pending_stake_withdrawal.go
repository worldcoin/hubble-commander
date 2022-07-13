package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
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
	return s.database.Badger.Insert(stake.BatchID, *stake)
}

func (s *PendingStakeWithdrawalStorage) RemovePendingStakeWithdrawal(batchID models.Uint256) error {
	var stake models.PendingStakeWithdrawal
	err := s.database.Badger.Delete(batchID, &stake)
	if errors.Is(err, bh.ErrNotFound) {
		return errors.WithStack(NewNotFoundError("pending stake withdrawal"))
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *PendingStakeWithdrawalStorage) GetReadyStateWithdrawals(currentBlock uint32) ([]models.PendingStakeWithdrawal, error) {
	stakes := make([]models.PendingStakeWithdrawal, 0)
	var stake models.PendingStakeWithdrawal
	err := s.database.Badger.Iterator(models.PendingStakeWithdrawalPrefix, db.PrefetchIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			err := item.Value(stake.SetBytes)
			if err != nil {
				return false, err
			}
			if stake.FinalisationBlock <= currentBlock {
				stakes = append(stakes, stake)
				return false, nil
			}
			return true, nil
		})
	if err != nil && !errors.Is(err, db.ErrIteratorFinished) {
		return nil, err
	}
	return stakes, nil
}
