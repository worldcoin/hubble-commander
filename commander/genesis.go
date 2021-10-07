package commander

import (
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	ErrRegisterGenesisAccountTimeout = fmt.Errorf("timeout")
	ErrGenesisAccountsUniqueStateID  = fmt.Errorf("accounts must have unique state IDs")
)

func AssignStateIDs(accounts []models.RegisteredGenesisAccount) []models.PopulatedGenesisAccount {
	populatedAccounts := make([]models.PopulatedGenesisAccount, 0, len(accounts))
	for i := range accounts {
		account := accounts[i]

		if account.Balance.CmpN(0) > 0 {
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
	seenStateIDs := make(map[uint32]bool)
	for i := range accounts {
		account := &accounts[i]

		if seenStateIDs[account.StateID] {
			return errors.WithStack(ErrGenesisAccountsUniqueStateID)
		}
		seenStateIDs[account.StateID] = true

		leaf := &models.AccountLeaf{
			PubKeyID:  account.PubKeyID,
			PublicKey: account.PublicKey,
		}
		_, err := saveSyncedSingleAccount(storage.AccountTree, leaf)
		if err != nil {
			return err
		}

		_, err = storage.StateTree.Set(account.StateID, &models.UserState{
			PubKeyID: account.PubKeyID,
			TokenID:  models.MakeUint256(0),
			Balance:  account.Balance,
			Nonce:    models.MakeUint256(0),
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
	txHashToPubKey := make(map[common.Hash]models.PublicKey)
	for i := range accounts {
		tx, err := eth.RegisterAccount(opts, accountRegistry, accounts[i].PublicKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		txHashToPubKey[tx.Hash()] = *accounts[i].PublicKey
		txs = append(txs, *tx)
	}

	registeredAccounts := make([]models.RegisteredGenesisAccount, len(accounts))
	accountsRegistered := 0
	for {
		select {
		case event, ok := <-registrations:
			if !ok {
				return nil, errors.WithStack(eth.ErrAccountWatcherIsClosed)
			}
			for i := range txs {
				if event.Raw.TxHash == txs[i].Hash() {
					registeredAccounts[i] = models.RegisteredGenesisAccount{
						GenesisAccount: accounts[i],
						PublicKey:      txHashToPubKey[event.Raw.TxHash],
						PubKeyID:       uint32(event.PubkeyID.Uint64()),
					}
					accountsRegistered += 1
				}
			}
			if accountsRegistered >= len(accounts) {
				return registeredAccounts, nil
			}

		case <-time.After(deployer.ChainTimeout):
			return nil, errors.WithStack(ErrRegisterGenesisAccountTimeout)
		}
	}
}

func (c *Commander) addGenesisBatch() error {
	batchID := models.MakeUint256(0)
	batch, err := c.storage.GetBatch(batchID)
	if batch != nil {
		return nil
	}
	if err != nil && !st.IsNotFoundError(err) {
		return err
	}

	root, err := c.storage.StateTree.Root()
	if err != nil {
		return err
	}

	batch, err = c.client.GetBatch(&batchID)
	if err != nil {
		return err
	}
	batch.PrevStateRoot = root

	return c.storage.AddBatch(batch)
}
