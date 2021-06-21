package commander

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommanderTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd *Commander
}

func (s *CommanderTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommanderTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	err := postgres.RecreateDatabase(cfg.Postgres)
	s.NoError(err)
	s.cmd = NewCommander(cfg)
	s.cmd.cfg.Ethereum = nil
}

func (s *CommanderTestSuite) TestStartStop() {
	s.False(s.cmd.IsRunning())

	err := s.cmd.Start()
	s.NoError(err)

	s.True(s.cmd.IsRunning())

	err = s.cmd.Stop()
	s.NoError(err)

	s.False(s.cmd.IsRunning())
}

func (s *CommanderTestSuite) TestStartAndWait() {
	stopped := false

	go func() {
		err := s.cmd.StartAndWait()
		s.NoError(err)
		stopped = true
	}()
	s.Eventually(func() bool {
		return s.cmd.IsRunning()
	}, 15*time.Second, 100*time.Millisecond)

	err := s.cmd.Stop()
	s.NoError(err)
	s.Eventually(func() bool {
		return stopped
	}, 1*time.Second, 100*time.Millisecond)
}

func (s *CommanderTestSuite) TestStart_SetsCorrectSyncedBlock() {
	err := s.cmd.Start()
	s.NoError(err)

	s.Equal(s.cmd.client.ChainState.DeploymentBlock-1, s.cmd.client.ChainState.SyncedBlock)

	err = s.cmd.Stop()
	s.NoError(err)
}

func (s *CommanderTestSuite) TestVerifyCommitment_ValidCommitmentRoot() {
	storage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	defer func() {
		err = storage.Teardown()
		s.NoError(err)
	}()
	testClient, err := eth.NewTestClient()
	s.NoError(err)
	defer testClient.Close()

	err = PopulateGenesisAccounts(storage.Storage, testClient.ChainState.GenesisAccounts)
	s.NoError(err)

	err = verifyCommitmentRoot(storage.Storage, testClient.Client)
	s.NoError(err)
}

func (s *CommanderTestSuite) TestVerifyCommitment_InvalidCommitmentRoot() {
	storage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	defer func() {
		err = storage.Teardown()
		s.NoError(err)
	}()
	testClient, err := eth.NewTestClient()
	s.NoError(err)
	defer testClient.Close()

	testClient.ChainState.GenesisAccounts = append(testClient.ChainState.GenesisAccounts, []models.PopulatedGenesisAccount{
		{
			PublicKey: models.PublicKey{5, 6, 7},
			PubKeyID:  1,
			StateID:   1,
			Balance:   models.MakeUint256(500),
		},
	}...)
	err = PopulateGenesisAccounts(storage.Storage, testClient.ChainState.GenesisAccounts)
	s.NoError(err)

	err = verifyCommitmentRoot(storage.Storage, testClient.Client)
	s.NotNil(err)
	s.Equal(ErrInvalidCommitmentRoot.Error(), err.Error())
}

func TestCommanderTestSuite(t *testing.T) {
	suite.Run(t, new(CommanderTestSuite))
}
