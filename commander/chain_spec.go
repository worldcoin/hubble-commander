package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"gopkg.in/yaml.v2"
)

func GenerateChainSpec(chainState *models.ChainState) (*string, error) {
	chainSpec := newChainSpec(chainState)

	yamlChainSpec, err := yaml.Marshal(chainSpec)
	if err != nil {
		return nil, err
	}

	return ref.String(string(yamlChainSpec)), nil
}

func newChainSpec(chainState *models.ChainState) models.ChainSpec {
	return models.ChainSpec{
		ChainID:         chainState.ChainID,
		AccountRegistry: chainState.AccountRegistry,
		Rollup:          chainState.Rollup,
		GenesisAccounts: chainState.GenesisAccounts,
	}
}
