package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

var selectedCommitmentCols = []string{
	"commitment.commitment_id",
	"commitment.transactions",
	"commitment.fee_receiver",
	"commitment.combined_signature",
	"commitment.post_state_root",
}

func (s *Storage) AddCommitment(commitment *models.Commitment) (*int32, error) {
	res := make([]int32, 0, 1)
	err := s.Postgres.Query(
		s.QB.Insert("commitment").
			Values(
				squirrel.Expr("DEFAULT"),
				commitment.Type,
				commitment.Transactions,
				commitment.FeeReceiver,
				commitment.CombinedSignature,
				commitment.PostStateRoot,
				commitment.IncludedInBatch,
			).
			Suffix("RETURNING commitment_id"),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return ref.Int32(res[0]), nil
}

func (s *Storage) GetCommitment(id int32) (*models.Commitment, error) {
	res := make([]models.Commitment, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("commitment").
			Where(squirrel.Eq{"commitment_id": id}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("commitment")
	}
	return &res[0], nil
}

func (s *Storage) MarkCommitmentAsIncluded(commitmentID int32, batchID models.Uint256) error {
	res, err := s.Postgres.Query(
		s.QB.Update("commitment").
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

func (s *Storage) GetCommitmentsByBatchID(batchID models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments := make([]models.CommitmentWithTokenID, 0, 32)
	err := s.Postgres.Query(
		s.QB.Select(selectedCommitmentCols...).
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
		stateLeaf, err := s.GetStateLeaf(commitments[i].FeeReceiverStateID)
		if err != nil {
			return nil, err
		}
		commitments[i].TokenID = stateLeaf.TokenIndex
	}

	return commitments, nil
}

func (s *Storage) DeleteCommitmentsByBatchIDs(batchID ...models.Uint256) error {
	res, err := s.Postgres.Query(
		s.QB.Delete("commitment").
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
