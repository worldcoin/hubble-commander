package config

import (
	"github.com/Worldcoin/hubble-commander/models"
	log "github.com/sirupsen/logrus"
)

const DefaultBlocksToFinalise = uint32(7 * 24 * 60 * 4)

func GetDeployerConfig() *DeployerConfig {
	setupViper("deployer-config")

	return &DeployerConfig{
		Bootstrap: &DeployerBootstrapConfig{
			BlocksToFinalise: getUint32("bootstrap.blocks_to_finalise", DefaultBlocksToFinalise), // nolint:misspell
			GenesisAccounts:  getGenesisAccounts(),
		},
		Ethereum: getEthereumConfig(),
	}
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := getString("bootstrap.genesis_path", "./genesis.yaml")
	log.Info("Reading gensis config from ", filename)
	return readGenesisAccounts(filename)
}

func readGenesisAccounts(filename string) []models.GenesisAccount {
	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Panicf("error reading genesis file: %s", err.Error())
	}
	return genesisAccounts
}
