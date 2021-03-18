package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

func (s *Storage) GetChainState(chainID models.Uint256) (*models.ChainState, error) {
	res := make([]models.ChainState, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("chain_state").
			Where(squirrel.Eq{"chain_id": chainID}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}

func (s *Storage) SetChainState(chainState *models.ChainState) error {
	_, err := s.DB.ExecBuilder(
		s.QB.
			Insert("chain_state").
			Values(chainState.ChainID, chainState.AccountRegistry, chainState.Rollup),
	)
	return err
}
