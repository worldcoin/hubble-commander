package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
)

type ChainStateStorage struct {
	database          *Database
	latestBlockNumber uint32
	syncedBlock       *uint64
}

func (s *ChainStateStorage) GetChainState(chainID models.Uint256) (*models.ChainState, error) {
	res := make([]models.ChainState, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("*").
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

func (s *ChainStateStorage) SetChainState(chainState *models.ChainState) error {
	_, err := s.database.Postgres.Query(
		s.database.QB.
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
