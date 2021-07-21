package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown   func() error
	testClient *eth.TestClient
	cmd        *Commander
}

func (s *AccountsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AccountsTestSuite) SetupTest() {
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

func (s *AccountsTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *AccountsTestSuite) TestSyncAccounts() {
	registrations, unsubscribe, err := s.testClient.WatchRegistrations(&bind.WatchOpts{Start: nil})
	s.NoError(err)
	defer unsubscribe()

	publicKey := models.PublicKey{2, 3, 4}
	pubKeyID, err := s.testClient.RegisterAccount(&publicKey, registrations)
	s.NoError(err)

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	err = s.cmd.syncAccounts(0, *latestBlockNumber)
	s.NoError(err)

	accounts, err := s.cmd.storage.GetAccountLeaves(&publicKey)
	s.NoError(err)
	s.Len(accounts, 1)
	s.Equal(*pubKeyID, accounts[0].PubKeyID)
}

func (s *AccountsTestSuite) TestSyncBatchAccount() {
	registrations, unsubscribe, err := s.testClient.WatchBatchAccountRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	var publicKeys [16]models.PublicKey
	for i := range publicKeys {
		publicKeys[i] = models.PublicKey{1, 1, byte(i)}
	}

	pubKeyIDs, err := s.testClient.RegisterBatchAccount(publicKeys, registrations)
	s.NoError(err)

	for i := range pubKeyIDs {
		leaves, err := s.cmd.storage.GetAccountLeaves(&publicKeys[i])
		s.NoError(err)
		s.Len(leaves, 1)
		s.Equal(pubKeyIDs[i], leaves[0].PubKeyID)
	}
}

func TestAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(AccountsTestSuite))
}
