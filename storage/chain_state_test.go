package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	chainState = models.ChainState{
		ChainID:                        models.MakeUint256(1),
		AccountRegistry:                utils.RandomAddress(),
		AccountRegistryDeploymentBlock: 234,
		TokenRegistry:                  utils.RandomAddress(),
		DepositManager:                 utils.RandomAddress(),
		WithdrawManager:                utils.RandomAddress(),
		Rollup:                         utils.RandomAddress(),
		SyncedBlock:                    10,
	}
)

type ChainStateTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *ChainStateTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ChainStateTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *ChainStateTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ChainStateTestSuite) TestSetChainState_SetAndRetrieve() {
	err := s.storage.SetChainState(&chainState)
	s.NoError(err)

	actual, err := s.storage.GetChainState()
	s.NoError(err)
	s.Equal(chainState, *actual)
}

func (s *ChainStateTestSuite) TestGetChainState_NotFound() {
	_, err := s.storage.GetChainState()
	s.ErrorIs(err, NewNotFoundError("chain state"))
	s.True(IsNotFoundError(err))
}

func TestChainStateTestSuite(t *testing.T) {
	suite.Run(t, new(ChainStateTestSuite))
}
