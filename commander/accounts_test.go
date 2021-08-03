package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
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
		Storage: testStorage.Storage,
		client:  s.testClient.Client,
	}
}

func (s *AccountsTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *AccountsTestSuite) TestSyncAccounts() {
	accounts := s.registerBatchAccount()
	accounts = append(accounts, s.registerSingleAccount())

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	err = s.cmd.syncAccounts(0, *latestBlockNumber)
	s.NoError(err)

	s.validateAccountsAfterSync(accounts)
}

func (s *AccountsTestSuite) TestSyncSingleAccounts() {
	account := s.registerSingleAccount()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	newAccountsCount, err := s.cmd.syncSingleAccounts(0, *latestBlockNumber)
	s.NoError(err)
	s.Equal(ref.Int(1), newAccountsCount)

	accounts, err := s.cmd.Storage.AccountTree.Leaves(&account.PublicKey)
	s.NoError(err)
	s.Len(accounts, 1)
	s.Equal(account.PubKeyID, accounts[0].PubKeyID)
}

func (s *AccountsTestSuite) TestSyncBatchAccounts() {
	accounts := s.registerBatchAccount()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	newAccountsCount, err := s.cmd.syncBatchAccounts(0, *latestBlockNumber)
	s.NoError(err)
	s.Equal(ref.Int(16), newAccountsCount)

	s.validateAccountsAfterSync(accounts)
}

func (s *AccountsTestSuite) registerSingleAccount() models.AccountLeaf {
	registrations, unsubscribe, err := s.testClient.WatchRegistrations(&bind.WatchOpts{Start: nil})
	s.NoError(err)
	defer unsubscribe()

	publicKey := models.PublicKey{2, 3, 4}
	pubKeyID, err := s.testClient.RegisterAccount(&publicKey, registrations)
	s.NoError(err)
	return models.AccountLeaf{
		PubKeyID:  *pubKeyID,
		PublicKey: publicKey,
	}
}

func (s *AccountsTestSuite) registerBatchAccount() []models.AccountLeaf {
	registrations, unsubscribe, err := s.testClient.WatchBatchAccountRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	publicKeys := make([]models.PublicKey, 16)
	for i := range publicKeys {
		publicKeys[i] = models.PublicKey{1, 1, byte(i)}
	}

	pubKeyIDs, err := s.testClient.RegisterBatchAccount(publicKeys, registrations)
	s.NoError(err)

	accounts := make([]models.AccountLeaf, 16)
	for i := range accounts {
		accounts[i] = models.AccountLeaf{
			PubKeyID:  pubKeyIDs[i],
			PublicKey: publicKeys[i],
		}
	}
	return accounts
}

func (s *AccountsTestSuite) validateAccountsAfterSync(accounts []models.AccountLeaf) {
	for i := range accounts {
		leaves, err := s.cmd.Storage.AccountTree.Leaves(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(leaves, 1)
		s.Equal(accounts[i].PubKeyID, leaves[0].PubKeyID)
	}
}

func TestAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(AccountsTestSuite))
}
