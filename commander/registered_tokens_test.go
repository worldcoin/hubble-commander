package commander

import (
	"github.com/Worldcoin/hubble-commander/eth"
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

func (s *RegisteredTokensTestSuite) TestSyncAccounts() {
	accounts := s.registerSingleToken()
}

func (s *RegisteredTokensTestSuite) registerSingleToken() models.RegisteredToken {
	s.testClient.Wa
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

func (s *RegisteredTokensTestSuite) registerBatchTokens() []models.RegisteredToken {
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
