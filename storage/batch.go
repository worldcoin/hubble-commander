package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
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
	return s.database.Badger.Insert(batch.ID, *models.NewStoredBatchFromBatch(batch))
}

func (s *BatchStorage) GetBatch(batchID models.Uint256) (*models.Batch, error) {
	var storedBatch models.StoredBatch
	err := s.database.Badger.Get(batchID, &storedBatch)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("batch"))
	}
	if err != nil {
		return nil, err
	}
	return storedBatch.ToBatch(), nil
}

func (s *BatchStorage) UpdateBatch(batch *models.Batch) error {
	err := s.database.Badger.Update(batch.ID, *models.NewStoredBatchFromBatch(batch))
	if err == bh.ErrNotFound {
		return errors.WithStack(NewNotFoundError("batch"))
	}
	return err
}

func (s *BatchStorage) GetMinedBatch(batchID models.Uint256) (*models.Batch, error) {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	if batch.Hash == nil {
		return nil, errors.WithStack(NewNotFoundError("batch"))
	}
	return batch, nil
}

func (s *BatchStorage) GetBatchByHash(batchHash common.Hash) (*models.Batch, error) {
	var storedBatch models.StoredBatch
	err := s.database.Badger.FindOne(
		&storedBatch,
		bh.Where("Hash").Eq(batchHash).Index("Hash"),
	)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("batch"))
	}
	if err != nil {
		return nil, err
	}

	return storedBatch.ToBatch(), nil
}

func (s *BatchStorage) GetLatestSubmittedBatch() (*models.Batch, error) {
	batch, err := s.reverseIterateBatches(func(batch *models.StoredBatch) bool {
		return batch.Hash != nil
	})
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *BatchStorage) GetNextBatchID() (*models.Uint256, error) {
	batch, err := s.reverseIterateBatches(func(batch *models.StoredBatch) bool {
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
	batch, err := s.reverseIterateBatches(func(batch *models.StoredBatch) bool {
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

	storedBatches := make([]models.StoredBatch, 0, 32)
	err := s.database.Badger.Find(&storedBatches, query)
	if err != nil {
		return nil, err
	}

	res := make([]models.Batch, 0, len(storedBatches))
	for i := range storedBatches {
		res = append(res, *storedBatches[i].ToBatch())
	}
	return res, nil
}

// DeleteBatches uses for loop instead badgerhold.DeleteMatching because it's faster
func (s *BatchStorage) DeleteBatches(batchIDs ...models.Uint256) error {
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		batch := models.StoredBatch{}
		for i := range batchIDs {
			err := txDatabase.Badger.Delete(batchIDs[i], batch)
			if err == bh.ErrNotFound {
				return errors.WithStack(NewNotFoundError("batch"))
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BatchStorage) reverseIterateBatches(filter func(batch *models.StoredBatch) bool) (*models.Batch, error) {
	var storedBatch models.StoredBatch
	err := s.database.Badger.Iterator(models.StoredBatchPrefix, db.ReversePrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
		err := item.Value(func(v []byte) error {
			return db.Decode(v, &storedBatch)
		})
		if err != nil {
			return false, err
		}
		return filter(&storedBatch), nil
	})
	if err == db.ErrIteratorFinished {
		return nil, errors.WithStack(NewNotFoundError("batch"))
	}
	if err != nil {
		return nil, err
	}
	return storedBatch.ToBatch(), nil
}
