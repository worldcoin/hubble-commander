package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeployerTestSuite struct {
	*require.Assertions
	suite.Suite
	sim *Simulator
}

func (s *DeployerTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DeployerTestSuite) SetupTest() {
	sim, err := NewSimulator()
	s.NoError(err)
	s.sim = sim
}

func (s *DeployerTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *DeployerTestSuite) TestEncodeTransferZero() {

}

func TestDeployerTestSuite(t *testing.T) {
	suite.Run(t, new(DeployerTestSuite))
}
