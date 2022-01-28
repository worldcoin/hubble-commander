package commander

import (
	"os"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth/chain"
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
	err := db.PruneDatabase(cfg.Badger)
	s.NoError(err)
	blockchain, err := GetChainConnection(cfg.Ethereum)
	s.NoError(err)
	s.prepareContracts(cfg, blockchain)
	s.cmd = NewCommander(cfg, blockchain)
}

func (s *CommanderTestSuite) TearDownTest() {
	err := os.Remove(s.chainSpecFile)
	s.NoError(err)
}

func (s *CommanderTestSuite) TestStartStop() {
	s.False(s.cmd.isRunning())

	err := s.cmd.Start()
	s.NoError(err)

	s.True(s.cmd.isRunning())

	err = s.cmd.Stop()
	s.NoError(err)

	s.False(s.cmd.isRunning())
}

func (s *CommanderTestSuite) TestStartAndWait() {
	var startAndWaitReturnTime *time.Time

	go func() {
		err := s.cmd.StartAndWait()
		s.NoError(err)
		startAndWaitReturnTime = ref.Time(time.Now())
	}()
	s.Eventually(s.cmd.isRunning, 15*time.Second, 100*time.Millisecond, "Commander hasn't started on time")

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

	s.Equal(s.cmd.client.ChainState.AccountRegistryDeploymentBlock-1, s.cmd.client.ChainState.SyncedBlock)

	err = s.cmd.Stop()
	s.NoError(err)
}

func (s *CommanderTestSuite) prepareContracts(cfg *config.Config, blockchain chain.Connection) {
	deployerCfg := config.GetDeployerConfig()
	yamlChainSpec, err := Deploy(deployerCfg, blockchain)
	s.NoError(err)

	file, err := os.CreateTemp("", "chain_spec_commander_test")
	s.NoError(err)

	_, err = file.WriteString(*yamlChainSpec)
	s.NoError(err)

	s.chainSpecFile = file.Name()
	cfg.Bootstrap.ChainSpecPath = &s.chainSpecFile
}

func TestCommanderTestSuite(t *testing.T) {
	suite.Run(t, new(CommanderTestSuite))
}
