package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddBatch(batch *models.Batch) error {
	_, err := s.Postgres.Query(
		s.QB.Insert("batch").
			Values(
				batch.ID,
				batch.Type,
				batch.TransactionHash,
				batch.Hash,
				batch.FinalisationBlock,
				batch.AccountTreeRoot,
				batch.PrevStateRoot,
				batch.Time,
			),
	).Exec()

	return err
}

func (s *Storage) MarkBatchAsSubmitted(batch *models.Batch) error {
	res, err := s.Postgres.Query(
		s.QB.Update("batch").
			Where(squirrel.Eq{"batch_id": batch.ID}).
			Set("batch_hash", batch.Hash).
			Set("finalisation_block", batch.FinalisationBlock). // nolint:misspell
			Set("account_tree_root", batch.AccountTreeRoot),
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

func (s *Storage) GetMinedBatch(batchID models.Uint256) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.And{
				squirrel.Eq{"batch_id": batchID},
				squirrel.NotEq{"batch_hash": nil},
			}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetBatch(batchID models.Uint256) (*models.Batch, error) {
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

func (s *Storage) GetBatchByHash(batchHash common.Hash) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("batch").
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

func (s *Storage) GetNextBatchID() (*models.Uint256, error) {
	res := make([]models.Uint256, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("batch_id").
			From("batch").
			OrderBy("batch_id DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return models.NewUint256(1), nil
	}
	return res[0].AddN(1), nil
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
		From("batch").
		OrderBy("batch_id")

	if from != nil {
		qb = qb.Where(squirrel.GtOrEq{"batch_id": from})
	}

	if to != nil {
		qb = qb.Where(squirrel.LtOrEq{"batch_id": to})
	}

	res := make([]models.Batch, 0, 32)
	err := s.Postgres.Query(qb).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) DeleteBatches(batchIDs ...models.Uint256) error {
	res, err := s.Postgres.Query(
		s.QB.Delete("batch").
			Where(squirrel.Eq{"batch_id": batchIDs}),
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
