package commander

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
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

	s.Equal(s.cmd.client.ChainState.AccountRegistryDeploymentBlock-1, s.cmd.client.ChainState.SyncedBlock)

	err = s.cmd.Stop()
	s.NoError(err)
}

func (s *CommanderTestSuite) TestStart_CreatesMempool() {
	tx := s.addTransaction()

	err := s.cmd.Start()
	s.NoError(err)

	txs := s.cmd.txPool.Pool.GetExecutableTxs(txtype.Transfer)
	s.Len(txs, 1)
	s.Equal(tx, txs[0])

	err = s.cmd.Stop()
	s.NoError(err)
}

func (s *CommanderTestSuite) addTransaction() *models.Transfer {
	storage, err := st.NewStorage(s.cmd.cfg)
	s.NoError(err)

	defer func() {
		s.cmd.cfg.Bootstrap.Prune = false
		err = storage.Close()
		s.NoError(err)
	}()

	tx := testutils.NewTransfer(0, 1, 0, 10)
	err = storage.AddTransaction(tx)
	s.NoError(err)

	return tx
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
