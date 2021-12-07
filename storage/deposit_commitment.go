package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/pkg/errors"
)

func (s *CommitmentStorage) AddDepositCommitment(commitment *models.DepositCommitment) error {
	return s.database.Badger.Insert(commitment.ID, stored.MakeCommitmentFromDepositCommitment(commitment))
}

func (s *CommitmentStorage) GetDepositCommitment(id *models.CommitmentID) (*models.DepositCommitment, error) {
	commitment, err := s.getStoredCommitment(id)
	if err != nil {
		return nil, err
	}
	if commitment.Type != batchtype.Deposit {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	return commitment.ToDepositCommitment(), nil
}
