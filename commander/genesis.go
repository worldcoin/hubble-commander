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
		registeredAccount, err := registerGenesisAccount(opts, accountRegistry, &accounts[i], ev)
		if err != nil {
			return nil, err
		}
		log.Printf("Registered genesis pubkey %s at %d", registeredAccount.PublicKey.String(), registeredAccount.PubKeyID)
		registeredAccounts = append(registeredAccounts, *registeredAccount)
	}

	return registeredAccounts, nil
}

func registerGenesisAccount(
	opts *bind.TransactOpts,
	accountRegistry *accountregistry.AccountRegistry,
	account *models.GenesisAccount,
	ev chan *accountregistry.AccountRegistryPubkeyRegistered,
) (*models.RegisteredGenesisAccount, error) {
	publicKey, err := bls.PrivateToPublicKey(account.PrivateKey)
	if err != nil {
		return nil, err
	}

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
				return &models.RegisteredGenesisAccount{
					GenesisAccount: *account,
					PublicKey:      *publicKey,
					PubKeyID:       uint32(event.PubkeyID.Uint64()),
				}, nil
			}
		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}
