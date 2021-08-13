package commander

import (
	"io/ioutil"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"gopkg.in/yaml.v2"
)

func ReadChainSpecFile(path string) (*models.ChainSpec, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var chainSpec models.ChainSpec
	err = yaml.Unmarshal(yamlFile, &chainSpec)
	if err != nil {
		return nil, err
	}

	return &chainSpec, nil
}

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
		DeploymentBlock: chainState.DeploymentBlock,
		Rollup:          chainState.Rollup,
		GenesisAccounts: chainState.GenesisAccounts,
	}
}

func makeChainStateFromChainSpec(chainSpec *models.ChainSpec) *models.ChainState {
	return &models.ChainState{
		ChainID:         chainSpec.ChainID,
		AccountRegistry: chainSpec.AccountRegistry,
		DeploymentBlock: chainSpec.DeploymentBlock,
		Rollup:          chainSpec.Rollup,
		GenesisAccounts: chainSpec.GenesisAccounts,
		SyncedBlock:     getInitialSyncedBlock(chainSpec.DeploymentBlock),
	}
}
