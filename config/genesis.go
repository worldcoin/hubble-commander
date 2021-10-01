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
		account := models.GenesisAccount{
			PublicKey:  nil,
			PrivateKey: nil,
			Balance:    models.MakeUint256(rawGenesisAccounts[i].Balance),
		}

		if rawGenesisAccounts[i].PublicKey != "" {
			decodedPublicKey, err := hex.DecodeString(rawGenesisAccounts[i].PublicKey)
			if err != nil {
				return nil, err
			}
			account.PublicKey = &models.PublicKey{}
			copy(account.PublicKey[:], decodedPublicKey)
		}

		if rawGenesisAccounts[i].PrivateKey != "" {
			decodedPrivateKey, err := hex.DecodeString(rawGenesisAccounts[i].PrivateKey)
			if err != nil {
				return nil, err
			}
			account.PrivateKey = &[32]byte{}
			copy(account.PrivateKey[:], decodedPrivateKey)
		}

		genesisAccounts = append(genesisAccounts, account)
	}

	return genesisAccounts, nil
}
