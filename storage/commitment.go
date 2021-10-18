package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
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

func (s *CommitmentStorage) AddCommitment(commitment *models.TxCommitment) error {
	return s.database.Badger.Insert(commitment.ID, models.MakeStoredCommitmentFromTxCommitment(commitment))
}

func (s *CommitmentStorage) GetCommitment(id *models.CommitmentID) (*models.TxCommitment, error) {
	commitment, err := s.GetStoredCommitment(id)
	if err != nil {
		return nil, err
	}
	return commitment.ToTxCommitment(), nil
}

func (s *Storage) GetCommitmentsByBatchID(batchID models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments, err := s.getStoredCommitmentsByBatchID(batchID)
	if err != nil {
		return nil, err
	}

	commitmentsWithToken := make([]models.CommitmentWithTokenID, 0, len(commitments))
	for i := range commitments {
		if commitments[i].Type != batchtype.Transfer && commitments[i].Type != batchtype.Create2Transfer {
			continue
		}
		commitment := commitments[i].ToTxCommitment()
		stateLeaf, err := s.StateTree.Leaf(commitment.FeeReceiver)
		if err != nil {
			return nil, err
		}
		commitmentsWithToken = append(commitmentsWithToken, models.CommitmentWithTokenID{
			ID:                 commitment.ID,
			PostStateRoot:      commitment.PostStateRoot,
			Transactions:       commitment.Transactions,
			FeeReceiverStateID: commitment.FeeReceiver,
			CombinedSignature:  commitment.CombinedSignature,
			TokenID:            stateLeaf.TokenID,
		})
	}

	return commitmentsWithToken, nil
}
