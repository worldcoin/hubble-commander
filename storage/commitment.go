package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bh "github.com/timshannon/badgerhold/v4"
)

func (s *CommitmentStorage) AddCommitment(commitment models.Commitment) error {
	switch commitment.GetCommitmentBase().Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return s.AddTxCommitment(commitment.ToTxCommitment())
	case batchtype.MassMigration:
		return s.AddMMCommitment(commitment.ToMMCommitment())
	case batchtype.Deposit:
		return s.addDepositCommitment(commitment.ToDepositCommitment())
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

func (s *CommitmentStorage) GetCommitmentsByBatchID(batchID models.Uint256) ([]models.Commitment, error) {
	storedCommitments, err := s.getStoredCommitmentsByBatchID(batchID)
	if err != nil {
		return nil, err
	}

	commitments := make([]models.Commitment, 0, len(storedCommitments))
	for i := range storedCommitments {
		switch storedCommitments[i].Type {
		case batchtype.Transfer, batchtype.Create2Transfer:
			commitments = append(commitments, storedCommitments[i].ToTxCommitment())
		case batchtype.MassMigration:
			commitments = append(commitments, storedCommitments[i].ToMMCommitment())
		case batchtype.Deposit:
			commitments = append(commitments, storedCommitments[i].ToDepositCommitment())
		default:
			panic("invalid commitment type")
		}
	}

	return commitments, nil
}

func (s *CommitmentStorage) UpdateCommitments(commitments []models.Commitment) error {
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		for i := range commitments {
			var commitment stored.Commitment

			switch commitments[i].GetCommitmentBase().Type {
			case batchtype.Transfer, batchtype.Create2Transfer:
				commitment = stored.MakeCommitmentFromTxCommitment(commitments[i].ToTxCommitment())
			case batchtype.MassMigration:
				commitment = stored.MakeCommitmentFromMMCommitment(commitments[i].ToMMCommitment())
			case batchtype.Deposit:
				commitment = stored.MakeCommitmentFromDepositCommitment(commitments[i].ToDepositCommitment())
			default:
				panic("invalid commitment type")
			}

			err := txDatabase.Badger.Update(commitments[i].GetCommitmentBase().ID, commitment)
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
