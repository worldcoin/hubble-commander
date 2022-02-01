package rollup

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils/consts"
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

func (s *RollupDeployerTestSuite) TestDeployConfiguredRollup_TransfersGenesisFunds() {
	deploymentCfg := DeploymentConfig{
		Params: Params{
			TotalGenesisAmount: models.NewUint256(5e9),
		},
	}
	rollupContracts, err := DeployConfiguredRollup(s.sim, &deploymentCfg)
	s.NoError(err)

	customToken, err := customtoken.NewTestCustomToken(rollupContracts.ExampleTokenAddress, s.sim.Backend)
	s.NoError(err)

	vaultAddress, err := rollupContracts.DepositManager.Vault(&bind.CallOpts{})
	s.NoError(err)

	vaultBalance, err := customToken.BalanceOf(&bind.CallOpts{}, vaultAddress)
	s.NoError(err)
	s.Equal(*deploymentCfg.TotalGenesisAmount.MulN(consts.L2Unit), models.MakeUint256FromBig(*vaultBalance))
}

func TestRollupDeployerTestSuite(t *testing.T) {
	suite.Run(t, new(RollupDeployerTestSuite))
}
