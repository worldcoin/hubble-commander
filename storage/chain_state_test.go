package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	chainState = models.ChainState{
		ChainID:         models.MakeUint256(1),
		AccountRegistry: common.HexToAddress("0x9f758331b439c1B664e86f2050F2360370F06849"),
		Rollup:          common.HexToAddress("0x1480c1b6bF90678820B259FCaFFbb751D3e3960B"),
		DeploymentBlock: 234,
		SyncedBlock:     10,
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
	s.storage, err = NewTestStorageWithoutPostgres()
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
	s.Equal(NewNotFoundError("chain state"), err)
	s.True(IsNotFoundError(err))
}

func TestChainStateTestSuite(t *testing.T) {
	suite.Run(t, new(ChainStateTestSuite))
}
