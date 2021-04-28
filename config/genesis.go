package config

import (
	"encoding/hex"
	"io/ioutil"
	"path"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"gopkg.in/yaml.v2"
)

func getGenesisAccounts(filename string) ([]models.GenesisAccount, error) {
	var rawGenesisAccount []models.RawGenesisAccount

	genesisFilePath := path.Join(utils.GetProjectRoot(), filename)

	yamlFile, err := ioutil.ReadFile(genesisFilePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &rawGenesisAccount)
	if err != nil {
		return nil, err
	}

	genesisAccounts := make([]models.GenesisAccount, 0, len(rawGenesisAccount))

	for i := range rawGenesisAccount {
		var privateKey [32]byte

		decodedPrivateKey, err := hex.DecodeString(rawGenesisAccount[i].PrivateKey)
		if err != nil {
			return nil, err
		}
		copy(privateKey[:], decodedPrivateKey)

		genesisAccounts = append(genesisAccounts, models.GenesisAccount{
			PrivateKey: privateKey,
			Balance:    rawGenesisAccount[i].Balance,
		})
	}

	return genesisAccounts, nil
}
