package commander

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrGenesisAccountsUniqueStateID = fmt.Errorf("accounts must have unique state IDs")

func AssignStateIDsAndCalculateTotalAmount(accounts []models.RegisteredGenesisAccount) (*models.Uint256, []models.PopulatedGenesisAccount) {
	totalGenesisAmount := models.NewUint256(0)
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
			totalGenesisAmount = totalGenesisAmount.Add(&account.Balance)
		}
	}
	return totalGenesisAmount, populatedAccounts
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

func RegisterGenesisAccounts(accountMgr *eth.AccountManager, accounts []models.GenesisAccount) ([]models.RegisteredGenesisAccount, error) {
	log.Println("Registering genesis accounts")
	txs := make([]types.Transaction, 0, len(accounts))
	for i := range accounts {
		tx, err := accountMgr.RegisterAccount(accounts[i].PublicKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		txs = append(txs, *tx)
	}

	receipts, err := chain.WaitForMultipleTxs(accountMgr.Blockchain.GetBackend(), txs...)
	if err != nil {
		return nil, err
	}

	registeredAccounts := make([]models.RegisteredGenesisAccount, 0, len(accounts))
	for i := range accounts {
		pubKeyID, err := accountMgr.RetrieveRegisteredPubKeyID(&receipts[i])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		registeredAccounts = append(registeredAccounts, models.RegisteredGenesisAccount{
			GenesisAccount: accounts[i],
			PublicKey:      *accounts[i].PublicKey,
			PubKeyID:       *pubKeyID,
		})
	}

	return registeredAccounts, nil
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
