package commander

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommanderTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd           *Commander
	chainSpecFile string
}

func (s *CommanderTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommanderTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	err := postgres.RecreateDatabase(cfg.Postgres)
	s.NoError(err)
	err = badger.PruneDatabase(cfg.Badger)
	s.NoError(err)
	chain, err := GetChainConnection(cfg.Ethereum)
	s.NoError(err)
	s.prepareContracts(cfg, chain)
	s.cmd = NewCommander(cfg, chain)
}

func (s *CommanderTestSuite) TearDownTest() {
	err := os.Remove(s.chainSpecFile)
	s.NoError(err)
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
	var startAndWaitReturnTime *time.Time

	go func() {
		err := s.cmd.StartAndWait()
		s.NoError(err)
		startAndWaitReturnTime = ref.Time(time.Now())
	}()
	s.Eventually(func() bool {
		return s.cmd.IsRunning()
	}, 15*time.Second, 100*time.Millisecond, "Commander hasn't started on time")

	err := s.cmd.Stop()
	s.NoError(err)
	stopReturnTime := time.Now()

	s.Eventually(func() bool {
		return startAndWaitReturnTime != nil
	}, 1*time.Second, 100*time.Millisecond, "StartAndWait hasn't returned on time")

	s.Greater(startAndWaitReturnTime.UnixNano(), stopReturnTime.UnixNano(), "Stop should return before StartAndWait")
}

func (s *CommanderTestSuite) TestStart_SetsCorrectSyncedBlock() {
	err := s.cmd.Start()
	s.NoError(err)

	s.Equal(s.cmd.client.ChainState.DeploymentBlock-1, s.cmd.client.ChainState.SyncedBlock)

	err = s.cmd.Stop()
	s.NoError(err)
}

func (s *CommanderTestSuite) prepareContracts(cfg *config.Config, chain deployer.ChainConnection) {
	yamlChainSpec, err := Deploy(cfg, chain)
	s.NoError(err)

	file, err := ioutil.TempFile("", "chain_spec_commander_test")
	s.NoError(err)

	_, err = file.Write([]byte(*yamlChainSpec))
	s.NoError(err)

	s.chainSpecFile = file.Name()
	cfg.Bootstrap.ChainSpecPath = &s.chainSpecFile
}

func TestCommanderTestSuite(t *testing.T) {
	suite.Run(t, new(CommanderTestSuite))
}
