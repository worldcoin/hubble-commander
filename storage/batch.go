package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddBatch(batch *models.Batch) error {
	_, err := s.DB.Query(
		s.QB.Insert("batch").
			Values(
				batch.Hash,
				batch.ID,
				batch.Type,
				batch.FinalisationBlock,
			),
	).Exec()
	return err
}

func (s *Storage) GetBatch(batchHash common.Hash) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.Eq{"batch_hash": batchHash}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetBatchByID(batchID models.Uint256) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.Eq{"batch_id": batchID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetBatchByCommitmentID(commitmentID int32) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.DB.Query(
		s.QB.Select("batch.*").
			From("batch").
			JoinClause("NATURAL JOIN commitment").
			Where(squirrel.Eq{"commitment_id": commitmentID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetLatestBatch() (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("batch").
			OrderBy("batch_id DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}

func (s *Storage) GetLatestFinalisedBatch(currentBlockNumber uint32) (*models.Batch, error) {
	res := make([]models.Batch, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("batch").
			Where(squirrel.LtOrEq{"finalisation_block": currentBlockNumber}).
			OrderBy("finalisation_block DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("batch")
	}
	return &res[0], nil
}
