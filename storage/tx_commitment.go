package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/pkg/errors"
)

func (s *CommitmentStorage) AddTxCommitment(commitment *models.TxCommitment) error {
	return s.database.Badger.Insert(commitment.ID, stored.MakeCommitmentFromTxCommitment(commitment))
}

func (s *CommitmentStorage) GetTxCommitment(id *models.CommitmentID) (*models.TxCommitment, error) {
	commitment, err := s.getStoredCommitment(id)
	if err != nil {
		return nil, err
	}
	if !s.isTxCommitmentType(commitment.Type) {
		return nil, errors.WithStack(NewNotFoundError("commitment"))
	}
	return commitment.ToTxCommitment(), nil
}

func (s *CommitmentStorage) GetTxCommitmentsByBatchID(batchID models.Uint256) ([]models.TxCommitment, error) {
	commitments, err := s.getStoredCommitmentsByBatchID(batchID)
	if err != nil {
		return nil, err
	}

	txCommitments := make([]models.TxCommitment, 0, len(commitments))
	for i := range commitments {
		if !s.isTxCommitmentType(commitments[i].Type) {
			continue
		}
		txCommitments = append(txCommitments, *commitments[i].ToTxCommitment())
	}

	return txCommitments, nil
}

func (s *CommitmentStorage) isTxCommitmentType(commitmentType batchtype.BatchType) bool {
	return commitmentType == batchtype.Transfer || commitmentType == batchtype.Create2Transfer
}
