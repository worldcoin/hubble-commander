package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

func (s *Storage) SetLatestBlockNumber(blockNumber uint32) {
	s.latestBlockNumber = blockNumber
}

func (s *Storage) GetLatestBlockNumber() uint32 {
	return s.latestBlockNumber
}

func (s *Storage) SetSyncedBlock(chainID models.Uint256, blockNumber uint32) error {
	s.syncedBlock = &blockNumber
	_, err := s.Postgres.Query(s.QB.Update("chain_state").
		Set("synced_block", blockNumber).
		Where(squirrel.Eq{"chain_id": chainID})).Exec()
	return err
}

func (s *Storage) GetSyncedBlock(chainID models.Uint256) (*uint32, error) {
	if s.syncedBlock != nil {
		return s.syncedBlock, nil
	}

	res := make([]uint32, 0, 1)
	err := s.Postgres.Query(s.QB.Select("synced_block").
		From("chain_state").
		Where(squirrel.Eq{"chain_id": chainID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return ref.Uint32(0), nil
	}
	return &res[0], nil
}
