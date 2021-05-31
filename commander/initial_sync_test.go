package commander

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InitialSyncTestSuite struct {
	*require.Assertions
	suite.Suite
	client   *eth.TestClient
	cmd      *Commander
	teardown func() error
}

func (s *InitialSyncTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *InitialSyncTestSuite) SetupTest() {
	storage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = storage.Teardown
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	cfg := config.GetTestConfig()
	cfg.Rollup.SyncSize = 1
	s.cmd = &Commander{
		cfg:     cfg,
		storage: storage.Storage,
		client:  s.client.Client,
	}
}

func (s *InitialSyncTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *InitialSyncTestSuite) TestInitialSync() {
	number, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	s.cmd.storage.SetLatestBlockNumber(*number + 3)

	accounts := []models.Account{
		{PublicKey: models.PublicKey{1, 1, 1}},
		{PublicKey: models.PublicKey{2, 2, 2}},
	}

	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()
	for i := range accounts {
		var senderPubKeyID *uint32
		senderPubKeyID, err = s.client.RegisterAccount(&accounts[i].PublicKey, registrations)
		s.NoError(err)
		s.Equal(uint32(i), *senderPubKeyID)
		accounts[i].PubKeyID = *senderPubKeyID
	}

	time.Sleep(200 * time.Millisecond)
	err = s.cmd.InitialSync()
	s.NoError(err)
	for i := range accounts {
		userAccounts, err := s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Contains(accounts, userAccounts[0])
	}
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(InitialSyncTestSuite))
}
