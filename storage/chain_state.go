package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
)

func (s *Storage) GetChainState(chainID models.Uint256) (*models.ChainState, error) {
	res := make([]models.ChainState, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("chain_state").
			Where(squirrel.Eq{"chain_id": chainID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("chain state")
	}
	return &res[0], nil
}

func (s *Storage) SetChainState(chainState *models.ChainState) error {
	_, err := s.Postgres.Query(
		s.QB.
			Insert("chain_state").
			Values(chainState.ChainID, chainState.AccountRegistry, chainState.Rollup),
	).Exec()
	return err
}

func (s *Storage) GetDomain(chainID models.Uint256) (*bls.Domain, error) {
	if s.domain != nil {
		return s.domain, nil
	}

	res := make([]bls.Domain, 0, 1)
	err := s.DB.Query(
		s.QB.Select("rollup").
			From("chain_state").
			Where(squirrel.Eq{"chain_id": chainID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("domain")
	}
	return &res[0], nil
}
