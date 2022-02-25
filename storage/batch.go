package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type BatchStorage struct {
	database *Database
}

func NewBatchStorage(database *Database) (*BatchStorage, error) {
	return &BatchStorage{
		database: database,
	}, nil
}

func (s *BatchStorage) copyWithNewDatabase(database *Database) *BatchStorage {
	newBatchStorage := *s
	newBatchStorage.database = database

	return &newBatchStorage
}

func (s *BatchStorage) AddBatch(batch *models.Batch) error {
	if batch.PrevStateRoot == nil {
		return ErrPrevStateRootIsEmpty
	}
	return s.database.Badger.Insert(batch.ID, *stored.NewBatchFromModelsBatch(batch))
}

func (s *BatchStorage) GetBatch(batchID models.Uint256) (*models.Batch, error) {
	var storedBatch stored.Batch
	err := s.database.Badger.Get(batchID, &storedBatch)
	if err == bh.ErrNotFound {
		return nil, errors.WithStack(NewNotFoundError("batch"))
	}
	if err != nil {
		return nil, err
	}
	return storedBatch.ToModelsBatch(), nil
}

func (s *BatchStorage) UpdateBatch(batch *models.Batch) error {
	if batch.PrevStateRoot == nil {
		return ErrPrevStateRootIsEmpty
	}
	err := s.database.Badger.Update(batch.ID, *stored.NewBatchFromModelsBatch(batch))
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
	storedBatches := make([]stored.Batch, 0, 1)
	// We're not using FindOne here because of inefficient underlying implementation
	err := s.database.Badger.Find(
		&storedBatches,
		bh.Where("Hash").Eq(&batchHash).Index("Hash").Limit(1),
	)
	if err != nil {
		return nil, err
	}

	if len(storedBatches) == 0 {
		return nil, errors.WithStack(NewNotFoundError("batch"))
	}

	return storedBatches[0].ToModelsBatch(), nil
}

func (s *BatchStorage) GetPendingBatches() ([]models.Batch, error) {
	var nilHash *common.Hash
	storedBatches := make([]stored.Batch, 0)
	err := s.database.Badger.Find(
		&storedBatches,
		bh.Where("Hash").Eq(nilHash).Index("Hash"),
	)
	if err != nil {
		return nil, err
	}

	batches := make([]models.Batch, 0, len(storedBatches))
	for i := range storedBatches {
		batches = append(batches, *storedBatches[i].ToModelsBatch())
	}
	return batches, nil
}

func (s *BatchStorage) GetLatestSubmittedBatch() (*models.Batch, error) {
	batch, err := s.reverseIterateBatches(func(batch *stored.Batch) bool {
		return batch.Hash != nil
	})
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *BatchStorage) GetNextBatchID() (*models.Uint256, error) {
	batch, err := s.reverseIterateBatches(func(batch *stored.Batch) bool {
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
	batch, err := s.reverseIterateBatches(func(batch *stored.Batch) bool {
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

	storedBatches := make([]stored.Batch, 0, 32)
	err := s.database.Badger.Find(&storedBatches, query)
	if err != nil {
		return nil, err
	}

	res := make([]models.Batch, 0, len(storedBatches))
	for i := range storedBatches {
		res = append(res, *storedBatches[i].ToModelsBatch())
	}
	return res, nil
}

// DeleteBatches uses for loop instead badgerhold.DeleteMatching because it's faster
func (s *BatchStorage) DeleteBatches(batchIDs ...models.Uint256) error {
	return s.database.ExecuteInTransaction(TxOptions{}, func(txDatabase *Database) error {
		storedBatch := stored.Batch{}
		for i := range batchIDs {
			err := txDatabase.Badger.Delete(batchIDs[i], storedBatch)
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

func (s *BatchStorage) reverseIterateBatches(filter func(batch *stored.Batch) bool) (*models.Batch, error) {
	var storedBatch stored.Batch
	err := s.database.Badger.Iterator(stored.BatchPrefix, db.ReversePrefetchIteratorOpts, func(item *bdg.Item) (bool, error) {
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
	return storedBatch.ToModelsBatch(), nil
}
