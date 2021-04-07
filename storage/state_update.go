package storage

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddStateUpdate(update *models.StateUpdate) error {
	_, err := s.DB.Query(
		s.QB.Insert("state_update").
			Columns(
				"merkle_path",
				"current_hash",
				"current_root",
				"prev_hash",
				"prev_root",
			).
			Values(
				update.MerklePath,
				update.CurrentHash,
				update.CurrentRoot,
				update.PrevHash,
				update.PrevRoot,
			),
	).Exec()

	return err
}

func (s *Storage) GetStateUpdateByRootHash(stateRootHash common.Hash) (*models.StateUpdate, error) {
	res := make([]models.StateUpdate, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("state_update").
			Where(squirrel.Eq{"current_root": stateRootHash}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("state update not found")
	}
	return &res[0], nil
}

func (s *Storage) DeleteStateUpdate(id uint64) error {
	_, err := s.DB.Query(
		s.QB.Delete("state_update").
			Where(squirrel.Eq{"id": id}),
	).Exec()

	return err
}
