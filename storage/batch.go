package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

type BatchStorage struct {
	database *Database
}

func NewBatchStorage(database *Database) *BatchStorage {
	return &BatchStorage{
		database: database,
	}
}

func (s *BatchStorage) copyWithNewDatabase(database *Database) *BatchStorage {
	newBatchStorage := *s
	newBatchStorage.database = database

	return &newBatchStorage
}

func (s *BatchStorage) AddBatch(batch *models.Batch) error {
	return s.database.Badger.Insert(batch.ID, *batch)
}

func (s *BatchStorage) GetBatch(batchID models.Uint256) (*models.Batch, error) {
	var batch models.Batch
	err := s.database.Badger.Get(batchID, &batch)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("batch")
	}
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (s *BatchStorage) MarkBatchAsSubmitted(batch *models.Batch) error {
	//TODO-bat: check if it's safe to upsert here
	return s.database.Badger.Upsert(batch.ID, *batch)
}

func (s *BatchStorage) GetMinedBatch(batchID models.Uint256) (*models.Batch, error) {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	if batch.Hash == nil {
		return nil, NewNotFoundError("batch")
	}
	return batch, nil
}

func (s *BatchStorage) GetBatchByHash(batchHash common.Hash) (*models.Batch, error) {
	batches := make([]models.Batch, 0, 1)
	err := s.database.Badger.Find(
		&batches,
		bh.Where("Hash").Eq(batchHash).Index("Hash"),
	)
	if err != nil {
		return nil, err
	}
	if len(batches) == 0 {
		return nil, NewNotFoundError("batch")
	}

	return &batches[0], nil
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
	criteria := bh.Where(bh.Key)
	var query *bh.Query
	if from == nil && to == nil {
		query = criteria.Ge(models.MakeUint256(0))
	}
	if from != nil && to != nil {
		query = criteria.Ge(*from).And(bh.Key).Le(*to)
	} else if from != nil {
		query = criteria.Ge(*from)
	} else if to != nil {
		query = criteria.Le(*to)
	}

	res := make([]models.Batch, 0, 32)
	err := s.database.Badger.Find(&res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *BatchStorage) DeleteBatches(batchIDs ...models.Uint256) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	batch := models.Batch{}
	for i := range batchIDs {
		err = txDatabase.Badger.Delete(batchIDs[i], batch)
		if err == bh.ErrNotFound {
			return NewNotFoundError("batch")
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
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
