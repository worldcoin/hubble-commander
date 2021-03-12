package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

type GenesisAccount struct {
	AccountIndex uint32 // TODO: Replace with pubkey
	Balance      models.Uint256
}

func PopulateGenesisAccounts(stateTree *storage.StateTree, accounts []GenesisAccount) error {
	for i, account := range accounts {
		err := stateTree.Set(uint32(i), &models.UserState{
			AccountIndex: models.MakeUint256(int64(account.AccountIndex)),
			TokenIndex:   models.MakeUint256(0),
			Balance:      account.Balance,
			Nonce:        models.MakeUint256(0),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
