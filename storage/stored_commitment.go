package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
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

func (s *CommitmentStorage) getStoredCommitment(id *models.CommitmentID) (*models.StoredCommitment, error) {
	commitment := new(models.StoredCommitment)
	err := s.database.Badger.Get(*id, commitment)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	if err != nil {
		return nil, err
	}
	return commitment, nil
}

func (s *CommitmentStorage) GetLatestCommitment() (*models.CommitmentBase, error) {
	var commitment *models.StoredCommitment
	var err error
	err = s.database.Badger.Iterator(models.StoredCommitmentPrefix, db.ReverseKeyIteratorOpts, func(item *bdg.Item) (bool, error) {
		commitment, err = decodeStoredCommitment(item)
		return true, err
	})
	if err == db.ErrIteratorFinished {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	if err != nil {
		return nil, err
	}

	return &commitment.CommitmentBase, nil
}

func (s *CommitmentStorage) DeleteCommitmentsByBatchIDs(batchIDs ...models.Uint256) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	var commitmentIDs []models.CommitmentID
	ids := make([]models.CommitmentID, 0, len(batchIDs))
	for i := range batchIDs {
		commitmentIDs, err = getCommitmentIDsByBatchID(txDatabase, db.ReverseKeyIteratorOpts, batchIDs[i])
		if err != nil {
			return err
		}
		ids = append(ids, commitmentIDs...)
	}

	if len(ids) == 0 {
		return errors.WithStack(NewNotFoundError("commitments"))
	}

	var commitment models.StoredCommitment
	for i := range ids {
		err = txDatabase.Badger.Delete(ids[i], commitment)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *CommitmentStorage) getStoredCommitmentsByBatchID(batchID models.Uint256) ([]models.StoredCommitment, error) {
	commitments := make([]models.StoredCommitment, 0, 32)
	prefix := getCommitmentPrefixWithBatchID(&batchID)
	err := s.database.Badger.Iterator(prefix, bdg.DefaultIteratorOptions, func(item *bdg.Item) (bool, error) {
		commitment, err := decodeStoredCommitment(item)
		if err != nil {
			return false, err
		}
		commitments = append(commitments, *commitment)
		return false, nil
	})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}
	return commitments, nil
}

func decodeStoredCommitment(item *bdg.Item) (*models.StoredCommitment, error) {
	var commitment models.StoredCommitment
	err := item.Value(commitment.SetBytes)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

func getCommitmentIDsByBatchID(txn *Database, opts bdg.IteratorOptions, batchID models.Uint256) ([]models.CommitmentID, error) {
	ids := make([]models.CommitmentID, 0, 32)
	prefix := getCommitmentPrefixWithBatchID(&batchID)
	err := txn.Badger.Iterator(prefix, opts, func(item *bdg.Item) (bool, error) {
		var id models.CommitmentID
		err := db.DecodeKey(item.Key(), &id, models.StoredCommitmentPrefix)
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
	commitmentPrefixLen := len(models.StoredCommitmentPrefix)
	prefix := make([]byte, commitmentPrefixLen+32)
	copy(prefix[:commitmentPrefixLen], models.StoredCommitmentPrefix)
	copy(prefix[commitmentPrefixLen:], batchID.Bytes())
	return prefix
}
