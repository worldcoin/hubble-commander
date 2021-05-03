package config

import (
	"encoding/hex"
	"io/ioutil"

	"github.com/Worldcoin/hubble-commander/models"
	"gopkg.in/yaml.v2"
)

func readGenesisFile(filepath string) ([]models.GenesisAccount, error) {
	var rawGenesisAccounts []models.RawGenesisAccount

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &rawGenesisAccounts)
	if err != nil {
		return nil, err
	}

	return decodeRawGenesisAccounts(rawGenesisAccounts)
}

func decodeRawGenesisAccounts(rawGenesisAccounts []models.RawGenesisAccount) ([]models.GenesisAccount, error) {
	genesisAccounts := make([]models.GenesisAccount, 0, len(rawGenesisAccounts))

	for i := range rawGenesisAccounts {
		var privateKey [32]byte

		decodedPrivateKey, err := hex.DecodeString(rawGenesisAccounts[i].PrivateKey)
		if err != nil {
			return nil, err
		}
		copy(privateKey[:], decodedPrivateKey)

		genesisAccounts = append(genesisAccounts, models.GenesisAccount{
			PrivateKey: privateKey,
			Balance:    rawGenesisAccounts[i].Balance,
		})
	}

	return genesisAccounts, nil
}
