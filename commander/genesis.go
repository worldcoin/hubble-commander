package commander

import (
	"fmt"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

type GenesisAccount struct {
	PublicKey models.PublicKey
	Balance   models.Uint256
}

type RegisteredGenesisAccount struct {
	GenesisAccount
	PubkeyID uint32
}

func PopulateGenesisAccounts(stateTree *storage.StateTree, accounts []RegisteredGenesisAccount) error {
	for i := range accounts {
		account := accounts[i]
		err := stateTree.Set(uint32(i), &models.UserState{
			PubkeyID:   account.PubkeyID,
			TokenIndex: models.MakeUint256(0),
			Balance:    account.Balance,
			Nonce:      models.MakeUint256(0),
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
		return nil, errors.WithStack(err)
	}
	defer sub.Unsubscribe()

	registeredAccounts := make([]RegisteredGenesisAccount, 0, len(accounts))

	for i := range accounts {
		registeredAccount, err := registerGenesisAccount(opts, accountRegistry, &accounts[i], ev)
		if err != nil {
			return nil, err
		}

		log.Printf("Registered genesis pubkey %s at %d", registeredAccount.PublicKey.String(), registeredAccount.PubkeyID)

		registeredAccounts = append(registeredAccounts, *registeredAccount)
	}

	return registeredAccounts, nil
}

func registerGenesisAccount(
	opts *bind.TransactOpts,
	accountRegistry *accountregistry.AccountRegistry,
	account *GenesisAccount,
	ev chan *accountregistry.AccountRegistryPubkeyRegistered,
) (*RegisteredGenesisAccount, error) {
	tx, err := accountRegistry.Register(opts, account.PublicKey.BigInts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for {
		select {
		case event, ok := <-ev:
			if !ok {
				return nil, errors.WithStack(fmt.Errorf("account event watcher is closed"))
			}
			if event.Raw.TxHash == tx.Hash() {
				pubkeyID := uint32(event.PubkeyID.Uint64())
				return &RegisteredGenesisAccount{
					GenesisAccount: *account,
					PubkeyID:       pubkeyID,
				}, nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}
