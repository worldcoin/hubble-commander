package commander

import (
	"math/big"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
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
	testStorage, err := st.NewTestStorageWithoutPostgres()
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
	registrations, unsubscribe, err := s.testClient.WatchTokenRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	tokenContract := s.testClient.CustomTokenAddress
	err = s.testClient.RequestRegisterToken(tokenContract)
	s.NoError(err)

	err = s.testClient.FinalizeRegisterToken(tokenContract)
	s.NoError(err)
	var tokenID *big.Int
Outer:
	for {

		select {
		case event, ok := <-registrations:
			if !ok {
				s.Fail("Token registry event watcher is closed")
			}
			if event.TokenContract == s.testClient.CustomTokenAddress {
				tokenID = event.TokenID
				break Outer
			}
		case <-time.After(deployer.ChainTimeout):
			s.Fail("Token registry event watcher timed out")
		}
	}
	return models.RegisteredToken{
		ID:       models.MakeUint256FromBig(*tokenID),
		Contract: tokenContract,
	}
}

func TestRegisteredTokensTestSuite(t *testing.T) {
	suite.Run(t, new(RegisteredTokensTestSuite))
}
