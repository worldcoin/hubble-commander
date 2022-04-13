package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Worldcoin/hubble-commander/models"
	"gopkg.in/yaml.v3"
)

type ModifiedUserState struct {
	models.UserState
	StateID uint32
}

func main() {
	accounts := readAccounts("accounts_dump.json")
	states := readUserStates("state_dump.json")

	fmt.Println("Creating genesis account items")

	var keysMissingState []uint32
	var genesisAccounts []models.GenesisAccount
	for i := range accounts {
		if i%1000 == 0 {
			fmt.Println(i)
		}

		state, err := findUserStateByPubKey(states, accounts[i].PubKeyID)
		if err != nil {
			keysMissingState = append(keysMissingState, accounts[i].PubKeyID)
		}
		genesisAccounts = append(genesisAccounts, models.GenesisAccount{
			PublicKey: accounts[i].PublicKey,
			StateID:   state.StateID,
			State: models.UserState{
				PubKeyID: accounts[i].PubKeyID,
				TokenID:  state.TokenID,
				Balance:  state.Balance,
				Nonce:    state.Nonce,
			}},
		)
	}

	fmt.Println("Writing missing keys file")
	missingKeysJSON, err := json.Marshal(keysMissingState)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("missing_keys.json", missingKeysJSON, 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing genesis account file")
	genesisAccountsYAML, err := yaml.Marshal(&genesisAccounts)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("gen_accounts.yaml", genesisAccountsYAML, 0600)
	if err != nil {
		panic(err)
	}
}

func findUserStateByPubKey(states []ModifiedUserState, pubKeyID uint32) (ModifiedUserState, error) {
	for _, v := range states {
		if v.PubKeyID == pubKeyID {
			return v, nil
		}
	}
	return ModifiedUserState{}, fmt.Errorf("not state found")
}

func readAccounts(file string) []models.AccountLeaf {
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	fileBytes, _ := io.ReadAll(jsonFile)

	var accounts []models.AccountLeaf
	err = json.Unmarshal(fileBytes, &accounts)
	if err != nil {
		panic(err)
	}

	return accounts
}

func readUserStates(file string) []ModifiedUserState {
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	fileBytes, _ := io.ReadAll(jsonFile)

	var accounts []ModifiedUserState
	err = json.Unmarshal(fileBytes, &accounts)
	if err != nil {
		panic(err)
	}

	return accounts
}
