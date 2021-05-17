package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

func PopulateGenesisAccounts(storage *st.Storage, accounts []models.RegisteredGenesisAccount) error {
	stateTree := st.NewStateTree(storage)

	for i := range accounts {
		account := accounts[i]
		err := storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  account.PubKeyID,
			PublicKey: account.PublicKey,
		})
		if err != nil {
			return err
		}

		if account.Balance.CmpN(0) == 1 {
			err = stateTree.Set(uint32(i), &models.UserState{
				PubKeyID:   account.PubKeyID,
				TokenIndex: models.MakeUint256(0),
				Balance:    account.Balance,
				Nonce:      models.MakeUint256(0),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RegisterGenesisAccounts(
	opts *bind.TransactOpts,
	accountRegistry *accountregistry.AccountRegistry,
	accounts []models.GenesisAccount,
) ([]models.RegisteredGenesisAccount, error) {
	registrations, unsubscribe, err := eth.WatchRegistrations(accountRegistry, nil)
	if err != nil {
		return nil, err
	}
	defer unsubscribe()

	registeredAccounts := make([]models.RegisteredGenesisAccount, 0, len(accounts))

	for i := range accounts {
		registeredAccount, err := registerGenesisAccount(opts, accountRegistry, &accounts[i], registrations)
		if err != nil {
			return nil, err
		}
		log.Printf("Registered genesis public key %s at id %d", registeredAccount.PublicKey.String(), registeredAccount.PubKeyID)
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

	pubKeyID, err := eth.RegisterAccount(opts, accountRegistry, publicKey, ev)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.RegisteredGenesisAccount{
		GenesisAccount: *account,
		PublicKey:      *publicKey,
		PubKeyID:       *pubKeyID,
	}, nil
}
