package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	for i, account := range accounts {
		if i%1000 == 0 {
			fmt.Println(i)
		}

		state, err := findUserStateByPubKey(states, account.PubKeyID)
		if err != nil {
			keysMissingState = append(keysMissingState, account.PubKeyID)
		}
		genesisAccounts = append(genesisAccounts, models.GenesisAccount{
			PublicKey: account.PublicKey,
			StateID:   state.StateID,
			State: models.UserState{
				PubKeyID: account.PubKeyID,
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
	err = ioutil.WriteFile("missing_keys.json", missingKeysJSON, 0666)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing genesis account file")
	genesisAccountsYAML, err := yaml.Marshal(&genesisAccounts)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("gen_accounts.yaml", genesisAccountsYAML, 0666)
	if err != nil {
		panic(err)
	}
}

func findUserStateByPubKey(states []ModifiedUserState, pubKeyId uint32) (ModifiedUserState, error) {
	for _, v := range states {
		if v.PubKeyID == pubKeyId {
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

	fileBytes, _ := ioutil.ReadAll(jsonFile)

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

	fileBytes, _ := ioutil.ReadAll(jsonFile)

	var accounts []ModifiedUserState
	err = json.Unmarshal(fileBytes, &accounts)
	if err != nil {
		panic(err)
	}

	return accounts
}
