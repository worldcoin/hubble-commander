package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (s *InternalStorage) GetChainState(chainID models.Uint256) (*models.ChainState, error) {
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

func (s *InternalStorage) SetChainState(chainState *models.ChainState) error {
	_, err := s.Postgres.Query(
		s.QB.
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

func (s *InternalStorage) GetDomain(chainID models.Uint256) (*bls.Domain, error) {
	if s.domain != nil {
		return s.domain, nil
	}

	res := make([]common.Address, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select("rollup as address").
			From("chain_state").
			Where(squirrel.Eq{"chain_id": chainID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("domain")
	}
	s.domain, err = bls.DomainFromBytes(crypto.Keccak256(res[0].Bytes()))
	return s.domain, err
}
