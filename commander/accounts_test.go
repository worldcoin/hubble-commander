package commander

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown   func() error
	storage    *st.Storage
	testClient *eth.TestClient
}

func (s *AccountsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AccountsTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
}

func (s *AccountsTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *AccountsTestSuite) TestWatchAccounts_PreviousAccounts() {
	publicKey := models.PublicKey{2, 3, 4}
	_, err := s.testClient.AccountRegistry.Register(s.testClient.Account, publicKey.BigInts())
	s.NoError(err)

	go func() {
		err = WatchAccounts(s.storage, s.testClient.Client, nil)
		s.NoError(err)
	}()

	var accounts []models.Account
	s.Eventually(func() bool {
		accounts, err = s.storage.GetAccounts(&publicKey)
		s.NoError(err)
		return len(accounts) == 1
	}, time.Second, testutils.TryInterval)
}

func (s *AccountsTestSuite) TestWatchAccounts_NewAccounts() {
	go func() {
		err := WatchAccounts(s.storage, s.testClient.Client, nil)
		s.NoError(err)
	}()

	time.Sleep(10 * time.Millisecond)

	publicKey := models.PublicKey{2, 3, 4}
	_, err := s.testClient.AccountRegistry.Register(s.testClient.Account, publicKey.BigInts())
	s.NoError(err)

	var accounts []models.Account
	s.Eventually(func() bool {
		accounts, err = s.storage.GetAccounts(&publicKey)
		s.NoError(err)
		return len(accounts) == 1
	}, time.Second, testutils.TryInterval)
}

func TestAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(AccountsTestSuite))
}
