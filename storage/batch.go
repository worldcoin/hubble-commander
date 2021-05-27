package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddPendingBatch(batch *models.PendingBatch) (*int32, error) {
	res := make([]int32, 0, 1)
	err := s.Postgres.Query(
		s.QB.Insert("batch").
			Values(
				squirrel.Expr("DEFAULT"),
				batch.Type,
				batch.TransactionHash,
			).
			Suffix("RETURNING batch_id"),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return ref.Int32(res[0]), nil
}

func (s *Storage) AddBatch(batch *models.Batch) (*int32, error) {
	res := make([]int32, 0, 1)
	err := s.Postgres.Query(
		s.QB.Insert("batch").
			Values(
				squirrel.Expr("DEFAULT"),
				batch.Type,
				batch.TransactionHash,
				batch.Hash,
				batch.Number,
				batch.FinalisationBlock,
			).
			Suffix("RETURNING batch_id"),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return ref.Int32(res[0]), nil
}

func (s *Storage) GetBatch(batchID int32) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.Eq{"batch_id": batchID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetBatchByNumber(batchNumber models.Uint256) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.Eq{"batch_number": batchNumber}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetBatchByCommitmentID(commitmentID int32) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("batch.*").
			From("batch").
			Join("commitment ON commitment.included_in_batch = batch.batch_id").
			Where(squirrel.Eq{"commitment_id": commitmentID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetLatestSubmittedBatch() (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.NotEq{"batch_hash": nil}).
			OrderBy("batch_id DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}

	return &res[0], nil
}

func (s *Storage) GetLatestFinalisedBatch(currentBlockNumber uint32) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.LtOrEq{"finalisation_block": currentBlockNumber}). // nolint:misspell
			OrderBy("finalisation_block DESC").                               // nolint:misspell
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetBatchesInRange(from, to *models.Uint256) ([]models.Batch, error) {
	qb := s.QB.Select("*").
		From("batch")

	if from != nil {
		qb = qb.Where(squirrel.GtOrEq{"batch_number": from})
	}

	if to != nil {
		qb = qb.Where(squirrel.LtOrEq{"batch_number": to})
	}

	res := make([]models.Batch, 0, 32)
	err := s.Postgres.Query(qb).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetBatchWithAccountRoot(batchHash common.Hash) (*models.BatchWithAccountRoot, error) {
	res := make([]models.BatchWithAccountRoot, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(
			"batch.*",
			"commitment.account_tree_root",
		).
			From("batch").
			Join("commitment ON commitment.included_in_batch = batch.batch_id").
			Where(squirrel.Eq{"batch_hash": batchHash}).
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetBatchWithAccountRootByNumber(batchNumber models.Uint256) (*models.BatchWithAccountRoot, error) {
	res := make([]models.BatchWithAccountRoot, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(
			"batch.*",
			"commitment.account_tree_root",
		).
			From("batch").
			Join("commitment ON commitment.included_in_batch = batch.batch_id").
			Where(squirrel.Eq{"batch_number": batchNumber}).
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}
