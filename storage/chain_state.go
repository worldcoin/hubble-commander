package storage

import (
	"sync/atomic"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type ChainStateStorage struct {
	database          *Database
	latestBlockNumber uint32
	syncedBlock       uint64
}

func NewChainStateStorage(database *Database) *ChainStateStorage {
	storage := &ChainStateStorage{database: database}
	atomic.StoreUint32(&storage.latestBlockNumber, 0)
	atomic.StoreUint64(&storage.syncedBlock, 0)

	return storage
}

func (s *ChainStateStorage) copyWithNewDatabase(database *Database) *ChainStateStorage {
	newChainStateStorage := &ChainStateStorage{database: database}
	atomic.StoreUint32(&newChainStateStorage.latestBlockNumber, s.GetLatestBlockNumber())
	atomic.StoreUint64(&newChainStateStorage.syncedBlock, 0)

	return newChainStateStorage
}

func (s *ChainStateStorage) GetChainState() (*models.ChainState, error) {
	var chainState models.ChainState
	err := s.database.Badger.Get("ChainState", &chainState)
	if errors.Is(err, bh.ErrNotFound) {
		return nil, errors.WithStack(NewNotFoundError("chain state"))
	}
	if err != nil {
		return nil, err
	}

	return &chainState, nil
}

func (s *ChainStateStorage) SetChainState(chainState *models.ChainState) error {
	return s.database.Badger.Upsert("ChainState", *chainState)
}
