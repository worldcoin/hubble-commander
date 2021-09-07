package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
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
	var batch models.Batch
	err := s.database.Badger.FindOne(
		&batch,
		bh.Where("Hash").Eq(batchHash).Index("Hash"),
	)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("batch")
	}
	if err != nil {
		return nil, err
	}

	return &batch, nil
}

func (s *BatchStorage) GetLatestSubmittedBatch() (*models.Batch, error) {
	batch, err := s.reverseIterateBatches(func(batch *models.Batch) bool {
		return batch.Hash != nil
	})
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *BatchStorage) GetNextBatchID() (*models.Uint256, error) {
	batch, err := s.reverseIterateBatches(func(batch *models.Batch) bool {
		return true
	})
	if IsNotFoundError(err) {
		return models.NewUint256(1), nil
	}
	if err != nil {
		return nil, err
	}
	return batch.ID.AddN(1), nil
}

func (s *BatchStorage) GetLatestFinalisedBatch(currentBlockNumber uint32) (*models.Batch, error) {
	batch, err := s.reverseIterateBatches(func(batch *models.Batch) bool {
		return batch.FinalisationBlock != nil && *batch.FinalisationBlock <= currentBlockNumber
	})
	if err != nil {
		return nil, err
	}
	return batch, nil
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

// DeleteBatches uses for loop instead badgerhold.DeleteMatching because it's faster
func (s *BatchStorage) DeleteBatches(batchIDs ...models.Uint256) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{})
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

func (s *BatchStorage) reverseIterateBatches(filter func(batch *models.Batch) bool) (*models.Batch, error) {
	var batch models.Batch
	err := s.database.Badger.Iterator(models.BatchPrefix, db.ReversePrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
		err := item.Value(func(v []byte) error {
			return db.Decode(v, &batch)
		})
		if err != nil {
			return false, err
		}
		return filter(&batch), nil
	})
	if err == db.ErrIteratorFinished {
		return nil, NewNotFoundError("batch")
	}
	if err != nil {
		return nil, err
	}
	return &batch, nil
}
