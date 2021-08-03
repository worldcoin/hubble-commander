package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type BatchStorage struct {
	database *Database
}

func (s *BatchStorage) BeginTransaction(opts TxOptions) (*db.TxController, *BatchStorage, error) {
	txController, txDatabase, err := s.database.BeginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txBatchStorage := *s
	txBatchStorage.database = txDatabase

	return txController, &txBatchStorage, nil
}

func (s *BatchStorage) AddBatch(batch *models.Batch) error {
	_, err := s.database.Postgres.Query(
		s.database.QB.Insert("batch").
			Values(
				batch.ID,
				batch.Type,
				batch.TransactionHash,
				batch.Hash,
				batch.FinalisationBlock,
				batch.AccountTreeRoot,
				batch.PrevStateRoot,
				batch.SubmissionTime,
			),
	).Exec()

	return err
}

func (s *BatchStorage) MarkBatchAsSubmitted(batch *models.Batch) error {
	res, err := s.database.Postgres.Query(
		s.database.QB.Update("batch").
			Where(squirrel.Eq{"batch_id": batch.ID}).
			Set("batch_hash", batch.Hash).
			Set("finalisation_block", batch.FinalisationBlock). // nolint:misspell
			Set("account_tree_root", batch.AccountTreeRoot).
			Set("submission_time", batch.SubmissionTime),
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

func (s *BatchStorage) GetMinedBatch(batchID models.Uint256) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("*").
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

func (s *BatchStorage) GetBatch(batchID models.Uint256) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("*").
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

func (s *BatchStorage) GetBatchByHash(batchHash common.Hash) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("*").
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

func (s *BatchStorage) GetLatestSubmittedBatch() (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("*").
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

func (s *BatchStorage) GetNextBatchID() (*models.Uint256, error) {
	res := make([]models.Uint256, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("batch_id").
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

func (s *BatchStorage) GetLatestFinalisedBatch(currentBlockNumber uint32) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("*").
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

func (s *BatchStorage) GetBatchesInRange(from, to *models.Uint256) ([]models.Batch, error) {
	qb := s.database.QB.Select("*").
		From("batch").
		OrderBy("batch_id")

	if from != nil {
		qb = qb.Where(squirrel.GtOrEq{"batch_id": from})
	}

	if to != nil {
		qb = qb.Where(squirrel.LtOrEq{"batch_id": to})
	}

	res := make([]models.Batch, 0, 32)
	err := s.database.Postgres.Query(qb).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *BatchStorage) DeleteBatches(batchIDs ...models.Uint256) error {
	res, err := s.database.Postgres.Query(
		s.database.QB.Delete("batch").
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

func (s *Storage) GetBatchByCommitmentID(commitmentID int32) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("batch.*").
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
