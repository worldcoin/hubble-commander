package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *CommitmentStorage) AddStoredCommitment(commitment *models.StoredCommitment) error {
	return s.database.Badger.Insert(commitment.ID, *commitment)
}

func (s *CommitmentStorage) GetStoredCommitment(id *models.CommitmentID) (*models.StoredCommitment, error) {
	commitment := &models.StoredCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: *id,
		},
	}
	err := s.database.Badger.Get(*id, commitment)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	if err != nil {
		return nil, err
	}
	return commitment, nil
}

func (s *CommitmentStorage) GetLatestCommitment() (*models.StoredCommitment, error) {
	var commitment *models.StoredCommitment
	var err error
	err = s.database.Badger.Iterator(models.StoredCommitmentPrefix, db.ReversePrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
		commitment, err = decodeStoredCommitment(item)
		return true, err
	})
	if err == db.ErrIteratorFinished {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	if err != nil {
		return nil, err
	}

	return commitment, nil
}

func decodeStoredCommitment(item *bdg.Item) (*models.StoredCommitment, error) {
	var commitment models.StoredCommitment
	err := item.Value(func(v []byte) error {
		return db.Decode(v, &commitment)
	})
	if err != nil {
		return nil, err
	}

	err = db.DecodeKey(item.Key(), &commitment.ID, models.StoredCommitmentPrefix)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}
