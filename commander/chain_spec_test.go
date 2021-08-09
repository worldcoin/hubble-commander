package commander

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type ChainSpecTestSuite struct {
	*require.Assertions
	suite.Suite
	storage    *st.TestStorage
	config     *cfg.Config
	chainState models.ChainState
	chainSpec  models.ChainSpec
}

func (s *ChainSpecTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ChainSpecTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.config = cfg.GetTestConfig()
	s.config.Ethereum = &cfg.EthereumConfig{
		ChainID: "1337",
	}
	s.chainState = models.ChainState{
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
}

func (s *ChainSpecTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ChainSpecTestSuite) TestGenerateChainSpec() {
	yamlChainSpec, err := GenerateChainSpec(&s.chainState)
	s.NoError(err)
	var chainSpec models.ChainSpec
	err = yaml.Unmarshal([]byte(*yamlChainSpec), &chainSpec)
	s.NoError(err)
	s.EqualValues(s.chainSpec, chainSpec)
}

func (s *ChainSpecTestSuite) TestLoadChainSpec() {
	chainState1337, err := s.storage.GetChainState(s.chainSpec.ChainID)
	s.NoError(err)

	s.chainSpec.DeploymentBlock = 483924
	s.chainSpec.ChainID = models.MakeUint256(4000)

	err = LoadChainSpec(s.config, &s.chainSpec)
	s.NoError(err)
	chainState4000, err := s.storage.GetChainState(s.chainSpec.ChainID)
	s.NoError(err)

	s.NotEqualValues(chainState1337, chainState4000)

	chainSpec := newChainSpec(chainState4000)
	s.EqualValues(s.chainSpec, chainSpec)
}

func (s *ChainSpecTestSuite) TestReadChainSpecFile() {
	yamlChainSpec, err := GenerateChainSpec(s.config)
	s.NoError(err)

	file, err := ioutil.TempFile("", "chain_state_test")
	s.NoError(err)
	defer func() {
		err = os.Remove(file.Name())
		s.NoError(err)
	}()

	_, err = file.Write([]byte(*yamlChainSpec))
	s.NoError(err)

	chainSpec, err := ReadChainSpecFile(file.Name())
	s.NoError(err)

	s.EqualValues(s.chainSpec, *chainSpec)
}

func TestChainSpecTestSuite(t *testing.T) {
	suite.Run(t, new(ChainSpecTestSuite))
}
