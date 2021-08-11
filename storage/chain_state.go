package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
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

func (s *ChainStateStorage) GetChainState(chainID models.Uint256) (*models.ChainState, error) {
	chainState := make([]models.ChainState, 0, 1)
	err := s.database.Badger.Find(
		&chainState,
		bh.Where("ChainID").Eq(chainID),
	)
	if err != nil {
		return nil, err
	}
	if len(chainState) == 0 {
		return nil, NewNotFoundError("chain state")
	}
	return &chainState[0], nil
}

func (s *ChainStateStorage) SetChainState(chainState *models.ChainState) error {
	return s.database.Badger.Upsert(chainState.ChainID, *chainState)
}
