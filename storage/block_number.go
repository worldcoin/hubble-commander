package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

func (s *StorageBase) SetLatestBlockNumber(blockNumber uint32) {
	s.latestBlockNumber = blockNumber
}

func (s *StorageBase) GetLatestBlockNumber() uint32 {
	return s.latestBlockNumber
}

func (s *StorageBase) SetSyncedBlock(chainID models.Uint256, blockNumber uint64) error {
	s.syncedBlock = &blockNumber
	_, err := s.Database.Postgres.Query(s.Database.QB.Update("chain_state").
		Set("synced_block", blockNumber).
		Where(squirrel.Eq{"chain_id": chainID})).Exec()
	return err
}

func (s *StorageBase) GetSyncedBlock(chainID models.Uint256) (*uint64, error) {
	if s.syncedBlock != nil {
		return s.syncedBlock, nil
	}

	res := make([]uint64, 0, 1)
	err := s.Database.Postgres.Query(s.Database.QB.Select("synced_block").
		From("chain_state").
		Where(squirrel.Eq{"chain_id": chainID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return ref.Uint64(0), nil
	}
	return &res[0], nil
}
