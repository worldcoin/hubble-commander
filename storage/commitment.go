package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (s *CommitmentStorage) AddCommitment(commitment models.Commitment) error {
	switch commitment.GetCommitmentBase().Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return s.AddTxCommitment(commitment.ToTxCommitment())
	case batchtype.MassMigration:
		return s.AddMMCommitment(commitment.ToMMCommitment())
	case batchtype.Deposit:
		return s.AddDepositCommitment(commitment.ToDepositCommitment())
	default:
		panic("invalid commitment type")
	}
}

func (s *CommitmentStorage) GetCommitment(id *models.CommitmentID) (models.Commitment, error) {
	commitment, err := s.getStoredCommitment(id)
	if err != nil {
		return nil, err
	}

	switch commitment.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return commitment.ToTxCommitment(), nil
	case batchtype.MassMigration:
		return commitment.ToMMCommitment(), nil
	case batchtype.Deposit:
		return commitment.ToDepositCommitment(), nil
	default:
		panic("invalid commitment type")
	}
}
