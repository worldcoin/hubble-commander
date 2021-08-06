package commander

import (
	"testing"

	cfg "github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type ChainSpecTestSuite struct {
	*require.Assertions
	suite.Suite
	storage   *st.TestStorage
	config    *cfg.Config
	chainSpec models.ChainSpec
}

func (s *ChainSpecTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ChainSpecTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.config = cfg.GetTestConfig()
	chainState := &models.ChainState{
		ChainID:         models.MakeUint256(1337),
		AccountRegistry: utils.RandomAddress(),
		DeploymentBlock: 8837,
		Rollup:          utils.RandomAddress(),
		GenesisAccounts: models.GenesisAccounts{
			{
				PublicKey: models.PublicKey{1, 2, 3, 4},
				PubKeyID:  17,
				StateID:   554,
				Balance:   models.MakeUint256(4534532),
			},
			{
				PublicKey: models.PublicKey{3, 4, 5, 6},
				PubKeyID:  93,
				StateID:   882,
				Balance:   models.MakeUint256(48391),
			},
			{
				PublicKey: models.PublicKey{5, 6, 7, 8},
				PubKeyID:  119,
				StateID:   1183,
				Balance:   models.MakeUint256(300920),
			},
		},
		SyncedBlock: 7738,
	}
	err = s.storage.SetChainState(chainState)
	s.NoError(err)
	s.chainSpec = newChainSpec(chainState)
	s.prepareConfig()
}

func (s *ChainSpecTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ChainSpecTestSuite) TestGenerateChainSpec() {
	yamlChainSpec, err := GenerateChainSpec(s.config)
	s.NoError(err)
	var chainSpec models.ChainSpec
	err = yaml.Unmarshal([]byte(*yamlChainSpec), &chainSpec)
	s.NoError(err)
	s.EqualValues(s.chainSpec, chainSpec)
}

func (s *ChainSpecTestSuite) prepareConfig() {
	config := *cfg.GetTestConfig()
	newEthConfig := *config.Ethereum
	newEthConfig.ChainID = "1337"
	config.Ethereum = &newEthConfig
	s.config = &config
}

func TestChainSpecTestSuite(t *testing.T) {
	suite.Run(t, new(ChainSpecTestSuite))
}
