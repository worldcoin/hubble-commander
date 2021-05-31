package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NewBlockTestSuite struct {
	*require.Assertions
	suite.Suite
	client   *eth.TestClient
	cmd      *Commander
	teardown func() error
}

func (s *NewBlockTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *NewBlockTestSuite) SetupTest() {
	var err error
	storage, err := st.NewTestStorage()
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

func (s *NewBlockTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *NewBlockTestSuite) TestSyncOnStart() {
	number, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	s.cmd.storage.SetLatestBlockNumber(*number + 3)

	// register sender account on chain
	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()
	senderPubKeyID, err := s.client.RegisterAccount(&models.PublicKey{1, 2, 3}, registrations)
	s.NoError(err)
	s.Equal(uint32(0), *senderPubKeyID)

	s.client.Commit()

	err = s.cmd.InitialSync()
	s.NoError(err)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(NewBlockTestSuite))
}
