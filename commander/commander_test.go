package commander

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/postgres"
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

func TestCommanderTestSuite(t *testing.T) {
	suite.Run(t, new(CommanderTestSuite))
}
