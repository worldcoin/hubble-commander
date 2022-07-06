package commander

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type ChainSpecTestSuite struct {
	*require.Assertions
	suite.Suite
	chainState *models.ChainState
	chainSpec  models.ChainSpec
}

func (s *ChainSpecTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ChainSpecTestSuite) SetupTest() {
	s.chainState = &models.ChainState{
		ChainID:                        models.MakeUint256(1337),
		AccountRegistry:                utils.RandomAddress(),
		AccountRegistryDeploymentBlock: 8837,
		TokenRegistry:                  utils.RandomAddress(),
		DepositManager:                 utils.RandomAddress(),
		WithdrawManager:                utils.RandomAddress(),
		Rollup:                         utils.RandomAddress(),
		GenesisAccounts: models.GenesisAccounts{
			{
				PublicKey: models.PublicKey{1, 2, 3, 4},
				StateID:   554,
				State: models.UserState{
					PubKeyID: 17,
					TokenID:  models.Uint256{},
					Balance:  models.MakeUint256(4534532),
					Nonce:    models.Uint256{},
				},
			},
			{
				PublicKey: models.PublicKey{3, 4, 5, 6},
				StateID:   882,
				State: models.UserState{
					PubKeyID: 93,
					TokenID:  models.Uint256{},
					Balance:  models.MakeUint256(48391),
					Nonce:    models.Uint256{},
				},
			},
			{
				PublicKey: models.PublicKey{5, 6, 7, 8},
				StateID:   1183,
				State: models.UserState{
					PubKeyID: 119,
					TokenID:  models.Uint256{},
					Balance:  models.MakeUint256(300920),
					Nonce:    models.Uint256{},
				},
			},
		},
	}
	s.chainSpec = makeChainSpec(s.chainState)
}

func (s *ChainSpecTestSuite) TestGenerateChainSpec() {
	yamlChainSpec, err := GenerateChainSpec(s.chainState)
	s.NoError(err)
	var chainSpec models.ChainSpec
	err = yaml.Unmarshal([]byte(*yamlChainSpec), &chainSpec)
	s.NoError(err)
	s.EqualValues(s.chainSpec, chainSpec)
}

func (s *ChainSpecTestSuite) TestReadChainSpecFile() {
	yamlChainSpec, err := GenerateChainSpec(s.chainState)
	s.NoError(err)

	file, err := os.CreateTemp("", "chain_spec_test")
	s.NoError(err)
	defer func() {
		err = os.Remove(file.Name())
		s.NoError(err)
	}()

	_, err = file.WriteString(*yamlChainSpec)
	s.NoError(err)

	chainSpec, err := ReadChainSpecFile(file.Name())
	s.NoError(err)

	s.EqualValues(s.chainSpec, *chainSpec)
}

func TestChainSpecTestSuite(t *testing.T) {
	suite.Run(t, new(ChainSpecTestSuite))
}
