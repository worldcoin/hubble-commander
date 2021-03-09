package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

func (s *Storage) AddStateUpdate(update *models.StateUpdate) error {
	_, err := s.DB.ExecBuilder(
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
	)

	return err
}

func (s *Storage) GetStateUpdate(id uint64) (*models.StateUpdate, error) {
	res := make([]models.StateUpdate, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("state_update").
			Where(squirrel.Eq{"id": id}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}
