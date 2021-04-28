package commander

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
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
	s.cmd = NewCommander(&cfg)
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
	time.Sleep(50 * time.Millisecond)

	err := s.cmd.Stop()
	s.NoError(err)
	s.Eventually(func() bool {
		return stopped
	}, 50*time.Millisecond, 10*time.Millisecond)
}

func TestCommanderTestSuite(t *testing.T) {
	suite.Run(t, new(CommanderTestSuite))
}
