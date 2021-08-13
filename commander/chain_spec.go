package commander

import (
	"io/ioutil"

	cfg "github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
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

func LoadChainSpec(config *cfg.Config, chainSpec *models.ChainSpec) error {
	chainState := newChainState(chainSpec)

	storage, err := st.NewStorage(config)
	if err != nil {
		return err
	}
	defer func() {
		err = storage.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = storage.GetChainState(chainState.ChainID)
	if err != nil {
		err = storage.SetChainState(&chainState)
		if err != nil {
			return err
		}
	}

	return nil
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

func newChainState(chainSpec *models.ChainSpec) models.ChainState {
	return models.ChainState{
		ChainID:         chainSpec.ChainID,
		AccountRegistry: chainSpec.AccountRegistry,
		DeploymentBlock: chainSpec.DeploymentBlock,
		Rollup:          chainSpec.Rollup,
		GenesisAccounts: chainSpec.GenesisAccounts,
		SyncedBlock:     getInitialSyncedBlock(chainSpec.DeploymentBlock),
	}
}
