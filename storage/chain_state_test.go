package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	chainState = models.ChainState{
		ChainID:                        models.MakeUint256(1),
		AccountRegistry:                common.HexToAddress("0x9f758331b439c1B664e86f2050F2360370F06849"),
		Rollup:                         common.HexToAddress("0x1480c1b6bF90678820B259FCaFFbb751D3e3960B"),
		AccountRegistryDeploymentBlock: 234,
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

	actual, err := s.storage.GetChainState(chainState.ChainID)
	s.NoError(err)
	s.Equal(chainState, *actual)
}

func (s *ChainStateTestSuite) TestGetChainState_NotFound() {
	_, err := s.storage.GetChainState(chainState.ChainID)
	s.Equal(NewNotFoundError("chain state"), err)
	s.True(IsNotFoundError(err))
}

func (s *ChainStateTestSuite) TestGetDomain() {
	err := s.storage.SetChainState(&chainState)
	s.NoError(err)

	expected, err := bls.DomainFromBytes(crypto.Keccak256(chainState.Rollup.Bytes()))
	s.NoError(err)

	domain, err := s.storage.GetDomain(chainState.ChainID)
	s.NoError(err)
	s.Equal(expected, domain)
	s.Equal(s.storage.domain, domain)
}

func (s *ChainStateTestSuite) TestGetDomain_NotFound() {
	domain, err := s.storage.GetDomain(chainState.ChainID)
	s.Equal(NewNotFoundError("domain"), err)
	s.Nil(domain)
}

func TestChainStateTestSuite(t *testing.T) {
	suite.Run(t, new(ChainStateTestSuite))
}
