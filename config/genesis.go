package config

import (
	"encoding/hex"
	"os"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"gopkg.in/yaml.v2"
)

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
		account := models.GenesisAccount{
			PublicKey:  nil,
			PrivateKey: nil,
			Balance:    models.MakeUint256(rawGenesisAccounts[i].Balance),
			State:      rawGenesisAccounts[i].State.ToStateLeaf(),
		}

		if rawGenesisAccounts[i].PublicKey != "" {
			decodedPublicKey, err := hex.DecodeString(rawGenesisAccounts[i].PublicKey)
			if err != nil {
				return nil, err
			}
			account.PublicKey = new(models.PublicKey)
			err = account.PublicKey.SetBytes(decodedPublicKey)
			if err != nil {
				return nil, err
			}
		}

		if rawGenesisAccounts[i].PrivateKey != "" {
			decodedPrivateKey, err := hex.DecodeString(rawGenesisAccounts[i].PrivateKey)
			if err != nil {
				return nil, err
			}
			account.PrivateKey = &[32]byte{}
			copy(account.PrivateKey[:], decodedPrivateKey)
		}

		if account.PrivateKey != nil && account.PublicKey != nil {
			if err := validateKeysMatch(*account.PrivateKey, account.PublicKey); err != nil {
				return nil, err
			}
		}

		genesisAccounts = append(genesisAccounts, account)
	}

	return genesisAccounts, nil
}

func validateKeysMatch(privateKey [32]byte, publicKey *models.PublicKey) error {
	derivedPublicKey, err := bls.PrivateToPublicKey(privateKey)
	if err != nil {
		return err
	}
	if *publicKey != *derivedPublicKey {
		return NewErrNonMatchingKeys(publicKey.String())
	}
	return nil
}
