package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RegisteredSpokesTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown   func() error
	testClient *eth.TestClient
	cmd        *Commander
}

func (s *RegisteredSpokesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RegisteredSpokesTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.cmd = &Commander{
		storage: testStorage.Storage,
		client:  s.testClient.Client,
		metrics: metrics.NewCommanderMetrics(),
	}
}

func (s *RegisteredSpokesTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *RegisteredSpokesTestSuite) TestSyncSingleSpoke() {
	registeredSpoke := s.registerSingleSpoke()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	err = s.cmd.syncSpokes(context.Background(), 0, *latestBlockNumber)
	s.NoError(err)

	syncedSpoke, err := s.cmd.storage.RegisteredSpokeStorage.GetRegisteredSpoke(registeredSpoke.ID)
	s.NoError(err)
	s.Equal(registeredSpoke.Contract, syncedSpoke.Contract)
	s.Equal(registeredSpoke.ID, syncedSpoke.ID)
}

func (s *RegisteredSpokesTestSuite) TestSyncSingleSpoke_CanSyncTheSameBlocksTwice() {
	_ = s.registerSingleSpoke()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	err = s.cmd.syncSpokes(context.Background(), 0, *latestBlockNumber)
	s.NoError(err)
	err = s.cmd.syncSpokes(context.Background(), 0, *latestBlockNumber)
	s.NoError(err)
}

func (s *RegisteredSpokesTestSuite) registerSingleSpoke() models.RegisteredSpoke {
	address := utils.RandomAddress()
	spokeID := RegisterSingleSpoke(s.Assertions, s.testClient, address)
	return models.RegisteredSpoke{
		ID:       *spokeID,
		Contract: address,
	}
}

func RegisterSingleSpoke(s *require.Assertions, testClient *eth.TestClient, spokeAddress common.Address) *models.Uint256 {
	spokeID, err := testClient.RegisterSpokeAndWait(spokeAddress)
	s.NoError(err)
	return spokeID
}

func TestRegisteredSpokesTestSuite(t *testing.T) {
	suite.Run(t, new(RegisteredSpokesTestSuite))
}
