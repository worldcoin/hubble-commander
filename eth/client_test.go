package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/testutils/deployer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	*require.Assertions
	suite.Suite
	sim    *simulator.Simulator
	client *Client
}

func (s *ClientTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ClientTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	contracts, err := deployer.DeployRollup(sim)
	s.NoError(err)
	s.client = NewClient(contracts.Rollup)
}

func (s *ClientTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *ClientTestSuite) Test_TODO() {
	// TODO
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
