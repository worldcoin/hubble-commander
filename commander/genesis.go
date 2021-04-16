package commander

import (
	"fmt"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

func PopulateGenesisAccounts(stateTree *storage.StateTree, accounts []models.RegisteredGenesisAccount) error {
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
	accounts []models.GenesisAccount,
) ([]models.RegisteredGenesisAccount, error) {
	ev := make(chan *accountregistry.AccountRegistryPubkeyRegistered)

	sub, err := accountRegistry.WatchPubkeyRegistered(&bind.WatchOpts{}, ev)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer sub.Unsubscribe()

	registeredAccounts := make([]models.RegisteredGenesisAccount, 0, len(accounts))

	for i := range accounts {
		wallet, err := bls.NewWallet(accounts[i].PrivateKey, bls.Domain{1, 2, 3})
		if err != nil {
			return nil, err
		}

		registeredAccount, err := registerGenesisAccount(opts, accountRegistry, &accounts[i], wallet.PublicKey(), ev)
		if err != nil {
			return nil, err
		}

		log.Printf("Registered genesis pubkey %s at %d", wallet.PublicKey().String(), registeredAccount.PubkeyID)

		registeredAccounts = append(registeredAccounts, *registeredAccount)
	}

	return registeredAccounts, nil
}

func registerGenesisAccount(
	opts *bind.TransactOpts,
	accountRegistry *accountregistry.AccountRegistry,
	account *models.GenesisAccount,
	publicKey *models.PublicKey,
	ev chan *accountregistry.AccountRegistryPubkeyRegistered,
) (*models.RegisteredGenesisAccount, error) {
	tx, err := accountRegistry.Register(opts, publicKey.BigInts())
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
				return &models.RegisteredGenesisAccount{
					GenesisAccount: *account,
					PubkeyID:       pubkeyID,
				}, nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}
