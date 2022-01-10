package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *CommitmentStorage) AddMMCommitment(commitment *models.MMCommitment) error {
	return s.database.Badger.Insert(commitment.ID, stored.MakeCommitmentFromMMCommitment(commitment))
}

func (s *CommitmentStorage) GetMMCommitment(id *models.CommitmentID) (*models.MMCommitment, error) {
	commitment, err := s.getStoredCommitment(id)
	if err != nil {
		return nil, err
	}
	if !s.isMMCommitmentType(commitment.Type) {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	return commitment.ToMMCommitment(), nil
}

func (s *CommitmentStorage) UpdateMMCommitments(commitments []models.MMCommitment) error {
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		for i := range commitments {
			err := txDatabase.Badger.Update(commitments[i].ID, stored.MakeCommitmentFromMMCommitment(&commitments[i]))
			if err == bh.ErrNotFound {
				return NewNotFoundError("commitment")
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *CommitmentStorage) GetMMCommitmentsByBatchID(batchID models.Uint256) ([]models.MMCommitment, error) {
	commitments, err := s.getStoredCommitmentsByBatchID(batchID)
	if err != nil {
		return nil, err
	}

	mmCommitments := make([]models.MMCommitment, 0, len(commitments))
	for i := range commitments {
		if !s.isMMCommitmentType(commitments[i].Type) {
			continue
		}
		mmCommitments = append(mmCommitments, *commitments[i].ToMMCommitment())
	}

	return mmCommitments, nil
}

func (s *CommitmentStorage) isMMCommitmentType(commitmentType batchtype.BatchType) bool {
	return commitmentType == batchtype.MassMigration
}
