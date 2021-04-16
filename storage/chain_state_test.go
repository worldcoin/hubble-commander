package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
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
	}
)

type ChainStateTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *ChainStateTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ChainStateTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *ChainStateTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *ChainStateTestSuite) Test_SetChainState() {
	err := s.storage.SetChainState(&chainState)
	s.NoError(err)

	actual, err := s.storage.GetChainState(chainState.ChainID)
	s.NoError(err)
	s.Equal(chainState, *actual)
}

func (s *ChainStateTestSuite) Test_GetChainState_NotFound() {
	_, err := s.storage.GetChainState(chainState.ChainID)
	s.Equal(NewNotFoundErr("chain state"), err)
	s.True(IsNotFoundErr(err))
}

func TestChainStateTestSuite(t *testing.T) {
	suite.Run(t, new(ChainStateTestSuite))
}
