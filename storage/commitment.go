package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

var selectedCommitmentCols = []string{
	"commitment.commitment_id",
	"commitment.transactions",
	"commitment.fee_receiver",
	"commitment.combined_signature",
	"commitment.post_state_root",
}

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

func (s *CommitmentStorage) AddCommitment(commitment *models.Commitment) error {
	err := s.database.Badger.Insert(models.CommitmentKey{
		BatchID:      commitment.BatchID,
		IndexInBatch: commitment.IndexInBatch,
	}, *commitment)
	return err
}

func (s *CommitmentStorage) GetCommitment(batchID models.Uint256, commitmentIndex uint32) (*models.Commitment, error) {
	var commitment models.Commitment
	err := s.database.Badger.Get(models.CommitmentKey{
		BatchID:      batchID,
		IndexInBatch: commitmentIndex,
	}, &commitment)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

func (s *CommitmentStorage) GetLatestCommitment() (*models.Commitment, error) {
	res := make([]models.Commitment, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("*").
			From("commitment").
			OrderBy("commitment_id DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("commitment")
	}
	return &res[0], nil
}

func (s *CommitmentStorage) MarkCommitmentAsIncluded(commitmentID int32, batchID models.Uint256) error {
	res, err := s.database.Postgres.Query(
		s.database.QB.Update("commitment").
			Where(squirrel.Eq{"commitment_id": commitmentID}).
			Set("included_in_batch", batchID),
	).Exec()
	if err != nil {
		return err
	}

	numUpdatedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numUpdatedRows == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *CommitmentStorage) DeleteCommitmentsByBatchIDs(batchID ...models.Uint256) error {
	res, err := s.database.Postgres.Query(
		s.database.QB.Delete("commitment").
			Where(squirrel.Eq{"included_in_batch": batchID}),
	).Exec()
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Storage) GetCommitmentsByBatchID(batchID models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments := make([]models.CommitmentWithTokenID, 0, 32)
	err := s.database.Postgres.Query(
		s.database.QB.Select(selectedCommitmentCols...).
			From("commitment").
			Where(squirrel.Eq{"included_in_batch": batchID}),
	).Into(&commitments)
	if err != nil {
		return nil, err
	}
	if len(commitments) == 0 {
		return nil, NewNotFoundError("commitments")
	}

	for i := range commitments {
		stateLeaf, err := s.StateTree.Leaf(commitments[i].FeeReceiverStateID)
		if err != nil {
			return nil, err
		}
		commitments[i].TokenID = stateLeaf.TokenID
	}

	return commitments, nil
}
