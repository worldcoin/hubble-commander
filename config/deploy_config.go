package config

import (
	"os"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	log "github.com/sirupsen/logrus"
)

const DefaultBlocksToFinalise = uint32(7 * 24 * 60 * 4)

func GetDeployerConfig() *DeployerConfig {
	setupViper("deployer-config")

	return &DeployerConfig{
		Bootstrap: &DeployerBootstrapConfig{
			BlocksToFinalise: getUint32("bootstrap.blocks_to_finalise", DefaultBlocksToFinalise), // nolint:misspell
			GenesisAccounts:  getGenesisAccounts(),
			Chooser:          getAddressOrNil("chooser_address"),
		},
		Ethereum: getEthereumConfig(),
	}
}

func GetDeployerTestConfig() *DeployerConfig {
	err := os.Chdir(utils.GetProjectRoot())
	if err != nil {
		panic(err)
	}

	return GetDeployerConfig()
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := getStringOrPanic("bootstrap.genesis_path")
	log.Printf("Reading genesis config from %s", filename)
	return readGenesisAccounts(filename)
}

func readGenesisAccounts(filename string) []models.GenesisAccount {
	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Panicf("error reading genesis file: %s", err.Error())
	}
	return genesisAccounts
}
