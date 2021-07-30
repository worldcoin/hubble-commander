package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

func (s *StorageBase) GetChainState(chainID models.Uint256) (*models.ChainState, error) {
	res := make([]models.ChainState, 0, 1)
	err := s.Database.Postgres.Query(
		s.Database.QB.Select("*").
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

func (s *StorageBase) SetChainState(chainState *models.ChainState) error {
	_, err := s.Database.Postgres.Query(
		s.Database.QB.
			Insert("chain_state").
			Values(
				chainState.ChainID,
				chainState.AccountRegistry,
				chainState.Rollup,
				chainState.GenesisAccounts,
				chainState.SyncedBlock,
				chainState.DeploymentBlock,
			),
	).Exec()
	return err
}
