package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

type ChainStateStorage struct {
	database          *Database
	latestBlockNumber uint32
	syncedBlock       *uint64
}

func NewChainStateStorage(database *Database) *ChainStateStorage {
	return &ChainStateStorage{
		database: database,
	}
}

func (s *ChainStateStorage) copyWithNewDatabase(database *Database) *ChainStateStorage {
	newChainStateStorage := *s
	newChainStateStorage.database = database

	return &newChainStateStorage
}

func (s *ChainStateStorage) GetChainState() (*models.ChainState, error) {
	var chainState models.ChainState
	err := s.database.Badger.Get("ChainState", &chainState)
	if err == bh.ErrNotFound {
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
