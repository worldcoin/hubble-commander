package commander

import (
	"fmt"
	"log"
	"time"

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
	for i := range accounts {
		account := accounts[i]
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

func RegisterGenesisAccounts(
	opts *bind.TransactOpts,
	accountRegistry *accountregistry.AccountRegistry,
	accounts []GenesisAccount,
) ([]RegisteredGenesisAccount, error) {
	ev := make(chan *accountregistry.AccountRegistryPubkeyRegistered)

	sub, err := accountRegistry.WatchPubkeyRegistered(&bind.WatchOpts{}, ev)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()

	registeredAccounts := make([]RegisteredGenesisAccount, 0, len(accounts))

	for _, account := range accounts {
		registeredAccount, err := registerGenesisAccount(opts, accountRegistry, account, ev)
		if err != nil {
			return nil, err
		}

		log.Printf("Registered genesis pubkey %s at %d", account.PublicKey.String(), registeredAccount.AccountIndex)

		registeredAccounts = append(registeredAccounts, *registeredAccount)
	}

	return registeredAccounts, nil
}

func registerGenesisAccount(opts *bind.TransactOpts, accountRegistry *accountregistry.AccountRegistry, account GenesisAccount, ev chan *accountregistry.AccountRegistryPubkeyRegistered) (*RegisteredGenesisAccount, error) {
	tx, err := accountRegistry.Register(opts, account.PublicKey.IntArray())
	if err != nil {
		return nil, err
	}

	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, fmt.Errorf("account event watcher is closed")
			}
			if event.Raw.TxHash == tx.Hash() {
				accountIndex := uint32(event.PubkeyID.Uint64())
				return &RegisteredGenesisAccount{
					GenesisAccount: account,
					AccountIndex:   accountIndex,
				}, nil
			}
		case <-time.After(500 * time.Millisecond):
			return nil, fmt.Errorf("timeout")
		}
	}
}
