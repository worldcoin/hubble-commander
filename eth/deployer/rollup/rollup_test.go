package rollup

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RollupDeployerTestSuite struct {
	*require.Assertions
	suite.Suite
	sim *simulator.Simulator
}

func (s *RollupDeployerTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RollupDeployerTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim
}

func (s *RollupDeployerTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *RollupDeployerTestSuite) TestDeployRollup() {
	rollupContracts, err := DeployRollup(s.sim)
	s.NoError(err)

	id, err := rollupContracts.Rollup.DomainSeparator(&bind.CallOpts{})
	s.NoError(err)

	var emptyBytes [32]byte
	s.NotEqual(emptyBytes, id)
}

func TestRollupDeployerTestSuite(t *testing.T) {
	suite.Run(t, new(RollupDeployerTestSuite))
}