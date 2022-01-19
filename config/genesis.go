package config

import (
	"os"

	"github.com/Worldcoin/hubble-commander/models"
	"gopkg.in/yaml.v2"
)

func readGenesisFile(filepath string) ([]models.GenesisAccount, error) {
	var genesisAccounts []models.GenesisAccount

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &genesisAccounts)
	if err != nil {
		return nil, err
	}
	return genesisAccounts, nil
}
