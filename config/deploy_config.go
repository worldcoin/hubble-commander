package config

import (
	"path"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	log "github.com/sirupsen/logrus"
)

const DefaultBlocksToFinalise = uint32(7 * 24 * 60 * 4)

func GetDeployerConfig() *DeployerConfig {
	setupViper(getDeployerConfigPath())

	return &DeployerConfig{
		Bootstrap: &DeployerBootstrapConfig{
			BlocksToFinalise: getUint32("bootstrap.blocks_to_finalise", DefaultBlocksToFinalise), // nolint:misspell
			GenesisAccounts:  getGenesisAccounts(),
		},
		Ethereum: getEthereumConfig(),
	}
}

func getDeployerConfigPath() string {
	return path.Join(utils.GetProjectRoot(), "deployer-config.yaml")
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := getString("bootstrap.genesis_path", getGenesisPath())
	return readGenesisAccounts(filename)
}

func readGenesisAccounts(filename string) []models.GenesisAccount {
	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Panicf("error reading genesis file: %s", err.Error())
	}
	return genesisAccounts
}

func getGenesisPath() string {
	return path.Join(utils.GetProjectRoot(), "genesis.yaml")
}
