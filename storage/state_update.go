package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

func (s *Storage) AddUpdate(update *models.StateUpdate) error {
	_, err := s.QB.Insert("state_update").
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
		).
		RunWith(s.DB).
		Exec()

	return err
}

func (s *Storage) GetUpdate(id uint64) (*models.StateUpdate, error) {
	res := make([]models.StateUpdate, 0, 1)
	sql, args, err := s.QB.Select("*").
		From("state_update").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}
