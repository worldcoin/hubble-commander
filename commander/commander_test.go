package commander

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
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
	cfg.Bootstrap.Prune = true
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
	s.False(s.cmd.isActive())

	err := s.cmd.Start()
	s.NoError(err)

	s.True(s.cmd.isActive())

	err = s.cmd.Stop()
	s.NoError(err)

	s.False(s.cmd.isActive())
}

func (s *CommanderTestSuite) TestStart_SetsCorrectSyncedBlock() {
	err := s.cmd.Start()
	s.NoError(err)

	blk, err := s.cmd.storage.GetSyncedBlock()
	s.NoError(err)

	s.Equal(s.cmd.client.ChainState.AccountRegistryDeploymentBlock-1, *blk)

	err = s.cmd.Stop()
	s.NoError(err)
}

func (s *CommanderTestSuite) prepareContracts(cfg *config.Config, blockchain chain.Connection) {
	deployerCfg := config.GetDeployerTestConfig()
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
