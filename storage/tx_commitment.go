package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/pkg/errors"
)

func (s *CommitmentStorage) AddTxCommitment(commitment *models.TxCommitment) error {
	return s.database.Badger.Insert(commitment.ID, models.MakeStoredCommitmentFromTxCommitment(commitment))
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

func (s *Storage) GetTxCommitmentsByBatchID(batchID models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments, err := s.getStoredCommitmentsByBatchID(batchID)
	if err != nil {
		return nil, err
	}

	commitmentsWithToken := make([]models.CommitmentWithTokenID, 0, len(commitments))
	for i := range commitments {
		if !s.isTxCommitmentType(commitments[i].Type) {
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

func (s *CommitmentStorage) isTxCommitmentType(commitmentType batchtype.BatchType) bool {
	return commitmentType == batchtype.Transfer || commitmentType == batchtype.Create2Transfer
}