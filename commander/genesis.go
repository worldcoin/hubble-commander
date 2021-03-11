package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

type GenesisAccount struct {
	accountIndex uint32 // TODO: Replace with pubkey
	balance      models.Uint256
}

func PopulateGenesisAccounts(stateTree *storage.StateTree, accounts []GenesisAccount) error {
	for i, account := range accounts {
		err := stateTree.Set(uint32(i), &models.UserState{
			AccountIndex: models.MakeUint256(int64(account.accountIndex)),
			TokenIndex:   models.MakeUint256(0),
			Balance:      account.balance,
			Nonce:        models.MakeUint256(0),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
