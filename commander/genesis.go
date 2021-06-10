package commander

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func AssignStateIDs(accounts []models.RegisteredGenesisAccount) []models.PopulatedGenesisAccount {
	populatedAccounts := make([]models.PopulatedGenesisAccount, 0, len(accounts))
	for i := range accounts {
		account := accounts[i]

		if account.Balance.CmpN(0) == 1 {
			populatedAccounts = append(populatedAccounts, models.PopulatedGenesisAccount{
				PublicKey: account.PublicKey,
				PubKeyID:  account.PubKeyID,
				StateID:   uint32(i),
				Balance:   account.Balance,
			})
		}
	}
	return populatedAccounts
}

func PopulateGenesisAccounts(storage *st.Storage, accounts []models.PopulatedGenesisAccount) error {
	stateTree := st.NewStateTree(storage)

	seenStateIDs := make(map[uint32]bool)
	for i := range accounts {
		account := &accounts[i]

		if seenStateIDs[account.StateID] {
			return errors.Errorf("accounts must have unique state IDs")
		}
		seenStateIDs[account.StateID] = true

		err := storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  account.PubKeyID,
			PublicKey: account.PublicKey,
		})
		if err != nil {
			return err
		}

		err = stateTree.Set(account.StateID, &models.UserState{
			PubKeyID:   account.PubKeyID,
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
	registrations, unsubscribe, err := eth.WatchRegistrations(accountRegistry, &bind.WatchOpts{})
	if err != nil {
		return nil, err
	}
	defer unsubscribe()

	txs := make([]types.Transaction, 0, len(accounts))
	for i := range accounts {
		account := &accounts[i]
		publicKey, err := bls.PrivateToPublicKey(account.PrivateKey)
		if err != nil {
			return nil, err
		}

		tx, err := eth.RegisterAccount(opts, accountRegistry, publicKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		txs = append(txs, *tx)
	}

	registeredAccounts := make([]models.RegisteredGenesisAccount, len(accounts))
	accountsRegistered := 0
	for {
		select {
		case event, ok := <-registrations:
			if !ok {
				return nil, errors.WithStack(fmt.Errorf("account event watcher is closed"))
			}
			for i := range txs {
				if event.Raw.TxHash == txs[i].Hash() {
					publicKey := models.MakePublicKeyFromInts(event.Pubkey)
					registeredAccounts[i] = models.RegisteredGenesisAccount{
						GenesisAccount: accounts[i],
						PublicKey:      publicKey,
						PubKeyID:       uint32(event.PubkeyID.Uint64()),
					}
					accountsRegistered += 1
				}
			}
			if accountsRegistered >= len(accounts) {
				return registeredAccounts, nil
			}

		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(fmt.Errorf("timeout"))
		}
	}
}
