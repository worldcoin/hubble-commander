package commander

import (
	"strconv"

	cfg "github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"gopkg.in/yaml.v2"
)

func GenerateChainSpec(config *cfg.Config) (*string, error) {
	storage, err := st.NewStorage(config)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = storage.Close()
		if err != nil {
			panic(err)
		}
	}()

	parsedChainID, err := strconv.ParseUint(config.Ethereum.ChainID, 10, 64)
	if err != nil {
		return nil, err
	}

	chainID := models.MakeUint256(parsedChainID)

	chainState, err := storage.GetChainState(chainID)
	if err != nil {
		return nil, err
	}

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
