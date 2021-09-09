package commander

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RegisteredTokensTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown   func() error
	testClient *eth.TestClient
	cmd        *Commander
}

func (s *RegisteredTokensTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RegisteredTokensTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.cmd = &Commander{
		storage: testStorage.Storage,
		client:  s.testClient.Client,
	}
}

func (s *RegisteredTokensTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *RegisteredTokensTestSuite) TestSyncSingleToken() {
	registeredToken := s.registerSingleToken()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	err = s.cmd.syncTokens(0, *latestBlockNumber)
	s.NoError(err)

	syncedToken, err := s.cmd.storage.RegisteredTokenStorage.GetRegisteredToken(models.MakeUint256(0))
	s.NoError(err)
	s.Equal(registeredToken.Contract, syncedToken.Contract)
	s.Equal(registeredToken.ID, syncedToken.ID)
}

func (s *RegisteredTokensTestSuite) registerSingleToken() models.RegisteredToken {
	token := models.RegisteredToken{
		Contract: s.testClient.ExampleTokenAddress,
	}
	RegisterSingleToken(s.Assertions, s.testClient, &token, ref.Uint64(0))
	return token
}

func RegisterSingleToken(s *require.Assertions, testClient *eth.TestClient, token *models.RegisteredToken, latestBlockNumber *uint64) {
	registrations, unsubscribe, err := testClient.WatchTokenRegistrations(&bind.WatchOpts{Start: latestBlockNumber})
	s.NoError(err)
	defer unsubscribe()

	err = testClient.RequestRegisterToken(token.Contract)
	s.NoError(err)

	err = testClient.FinalizeRegisterToken(token.Contract)
	s.NoError(err)
	for {
		select {
		case event, ok := <-registrations:
			if !ok {
				s.Fail("Token registry event watcher is closed")
			}
			if event.TokenContract == testClient.ExampleTokenAddress {
				token.ID = models.MakeUint256FromBig(*event.TokenID)
				return
			}
		case <-time.After(100 * time.Millisecond):
			s.Fail("Token registry event watcher timed out")
		}
	}
}

func TestRegisteredTokensTestSuite(t *testing.T) {
	suite.Run(t, new(RegisteredTokensTestSuite))
}
