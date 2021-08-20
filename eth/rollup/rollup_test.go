package rollup

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeployerTestSuite struct {
	*require.Assertions
	suite.Suite
	sim *simulator.Simulator
}

func (s *DeployerTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DeployerTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator(&config.EthereumConfig{})
	s.NoError(err)
	s.sim = sim
}

func (s *DeployerTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *DeployerTestSuite) TestNewRollup() {
	rollupContracts, err := DeployRollup(s.sim)
	s.NoError(err)

	id, err := rollupContracts.Rollup.DomainSeparator(&bind.CallOpts{})
	s.NoError(err)

	var emptyBytes [32]byte
	s.NotEqual(emptyBytes, id)
}

func TestDeployerTestSuite(t *testing.T) {
	suite.Run(t, new(DeployerTestSuite))
}
