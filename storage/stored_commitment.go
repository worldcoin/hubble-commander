package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type CommitmentStorage struct {
	database *Database
}

func NewCommitmentStorage(database *Database) *CommitmentStorage {
	return &CommitmentStorage{
		database: database,
	}
}

func (s *CommitmentStorage) copyWithNewDatabase(database *Database) *CommitmentStorage {
	newCommitmentStorage := *s
	newCommitmentStorage.database = database

	return &newCommitmentStorage
}

func (s *CommitmentStorage) getStoredCommitment(id *models.CommitmentID) (*stored.Commitment, error) {
	storedCommitment := new(stored.Commitment)
	err := s.database.Badger.Get(*id, storedCommitment)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	if err != nil {
		return nil, err
	}
	return storedCommitment, nil
}

func (s *CommitmentStorage) GetLatestCommitment() (*models.CommitmentBase, error) {
	var storedCommitment *stored.Commitment
	var err error
	err = s.database.Badger.Iterator(stored.CommitmentPrefix, db.ReverseKeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		storedCommitment, err = decodeStoredCommitment(item)
		return true, err
	})
	if err == db.ErrIteratorFinished {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	if err != nil {
		return nil, err
	}

	return &storedCommitment.CommitmentBase, nil
}

func (s *CommitmentStorage) DeleteCommitmentsByBatchIDs(batchIDs ...models.Uint256) error {
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		ids := make([]models.CommitmentID, 0, len(batchIDs))
		for i := range batchIDs {
			commitmentIDs, err := getCommitmentIDsByBatchID(txDatabase, batchIDs[i])
			if err != nil {
				return err
			}
			ids = append(ids, commitmentIDs...)
		}

		if len(ids) == 0 {
			return errors.WithStack(NewNotFoundError("commitments"))
		}

		var storedCommitment stored.Commitment
		for i := range ids {
			err := txDatabase.Badger.Delete(ids[i], storedCommitment)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *CommitmentStorage) getStoredCommitmentsByBatchID(batchID models.Uint256) ([]stored.Commitment, error) {
	storedCommitments := make([]stored.Commitment, 0, 32)
	prefix := getCommitmentPrefixWithBatchID(&batchID)
	err := s.database.Badger.Iterator(prefix, bdg.DefaultIteratorOptions, func(item *bdg.Item) (bool, error) {
		commitment, err := decodeStoredCommitment(item)
		if err != nil {
			return false, err
		}
		storedCommitments = append(storedCommitments, *commitment)
		return false, nil
	})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}
	return storedCommitments, nil
}

func decodeStoredCommitment(item *bdg.Item) (*stored.Commitment, error) {
	var storedCommitment stored.Commitment
	err := item.Value(storedCommitment.SetBytes)
	if err != nil {
		return nil, err
	}
	return &storedCommitment, nil
}

func getCommitmentIDsByBatchID(txn *Database, batchID models.Uint256) ([]models.CommitmentID, error) {
	ids := make([]models.CommitmentID, 0, 32)
	prefix := getCommitmentPrefixWithBatchID(&batchID)
	err := txn.Badger.Iterator(prefix, db.KeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		var id models.CommitmentID
		err := db.DecodeKey(item.Key(), &id, stored.CommitmentPrefix)
		if err != nil {
			return false, err
		}
		ids = append(ids, id)
		return false, nil
	})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}
	return ids, nil
}

func getCommitmentPrefixWithBatchID(batchID *models.Uint256) []byte {
	storedCommitmentPrefixLen := len(stored.CommitmentPrefix)
	prefix := make([]byte, storedCommitmentPrefixLen+32)
	copy(prefix[:storedCommitmentPrefixLen], stored.CommitmentPrefix)
	copy(prefix[storedCommitmentPrefixLen:], batchID.Bytes())
	return prefix
}

func (s *CommitmentStorage) GetCommitmentsByBatchID(batchID models.Uint256, commitmentType batchtype.BatchType) (
	[]models.Commitment,
	error,
) {
	if commitmentType == batchtype.Deposit {
		commitment, err := s.GetDepositCommitment(&models.CommitmentID{
			BatchID:      batchID,
			IndexInBatch: 0,
		})
		if err != nil {
			return nil, err
		}
		return []models.Commitment{commitment}, nil
	}
	txCommitments, err := s.GetTxCommitmentsByBatchID(batchID)
	if err != nil {
		return nil, err
	}
	commitments := make([]models.Commitment, 0, len(txCommitments))
	for i := range txCommitments {
		commitments = append(commitments, &txCommitments[i])
	}
	return commitments, nil
}
