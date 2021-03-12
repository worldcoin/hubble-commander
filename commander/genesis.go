package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type GenesisAccount struct {
	PublicKey models.PublicKey
	Balance   models.Uint256
}

type RegisteredGenesisAccount struct {
	GenesisAccount
	AccountIndex uint32
}

func PopulateGenesisAccounts(stateTree *storage.StateTree, accounts []RegisteredGenesisAccount) error {
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

func RegisterGenesisAccounts(opts *bind.TransactOpts, accountRegistry *accountregistry.AccountRegistry, accounts []GenesisAccount) ([]RegisteredGenesisAccount, error) {
	ev := make(chan *accountregistry.AccountRegistryPubkeyRegistered)

	sub, err := accountRegistry.AccountRegistryFilterer.WatchPubkeyRegistered(&bind.WatchOpts{}, ev)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	registeredAccounts := make([]RegisteredGenesisAccount, 0, len(accounts))

	for _, account := range accounts {
		tx, err := accountRegistry.Register(opts, account.PublicKey.IntArray())
		if err != nil {
			return nil, err
		}
		var accountIndex uint32
		for {
			event, ok := <-ev
			if !ok {
				log.Fatal("Account event watcher is closed")
			}
			if event.Raw.TxHash == tx.Hash() {
				accountIndex = uint32(event.PubkeyID.Uint64())
				break
			}
		}

		log.Printf("Registered genesis pubkey %s at %d", account.PublicKey.String(), accountIndex)

		registeredAccounts = append(registeredAccounts, RegisteredGenesisAccount{
			GenesisAccount: account,
			AccountIndex:   accountIndex,
		})
	}

	return registeredAccounts, nil
}
