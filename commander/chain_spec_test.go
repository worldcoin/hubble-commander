package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestChainSpecTestSuite(t *testing.T) {
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
	expectedChainSpec := newChainSpec(chainState)

	yamlChainSpec, err := GenerateChainSpec(chainState)
	require.NoError(t, err)
	var chainSpec models.ChainSpec
	err = yaml.Unmarshal([]byte(*yamlChainSpec), &chainSpec)
	require.NoError(t, err)
	require.EqualValues(t, expectedChainSpec, chainSpec)
}
