package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/pkg/errors"
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

func (s *CommitmentStorage) isMMCommitmentType(commitmentType batchtype.BatchType) bool {
	return commitmentType == batchtype.MassMigration
}
