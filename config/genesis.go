package config

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var errMissingGenesisPublicKey = fmt.Errorf("genesis accounts require public keys")

func readGenesisFile(filepath string) ([]models.GenesisAccount, error) {
	var rawGenesisAccounts []models.RawGenesisAccount

	yamlFile, err := os.ReadFile(filepath)
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
		if rawGenesisAccounts[i].PublicKey == "" {
			return nil, errors.WithStack(errMissingGenesisPublicKey)
		}

		decodedPublicKey, err := hex.DecodeString(rawGenesisAccounts[i].PublicKey)
		if err != nil {
			return nil, err
		}
		publicKey := models.PublicKey{}
		err = publicKey.SetBytes(decodedPublicKey)
		if err != nil {
			return nil, err
		}

		genesisAccounts = append(genesisAccounts, models.GenesisAccount{
			PublicKey: publicKey,
			State:     rawGenesisAccounts[i].State.ToStateLeaf(),
		})
	}

	return genesisAccounts, nil
}
