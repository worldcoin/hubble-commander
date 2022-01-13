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

func PopulateGenesisAccounts(storage *st.Storage, accounts []models.PopulatedGenesisAccount) error {
	seenStateIDs := make(map[uint32]bool)
	for i := range accounts {
		account := &accounts[i]

		if seenStateIDs[account.StateID] {
			return errors.WithStack(ErrGenesisAccountsUniqueStateID)
		}
		seenStateIDs[account.StateID] = true

		leaf := &models.AccountLeaf{
			PubKeyID:  account.State.PubKeyID,
			PublicKey: account.PublicKey,
		}
		_, err := saveSyncedSingleAccount(storage.AccountTree, leaf)
		if err != nil {
			return err
		}

		_, err = storage.StateTree.Set(account.StateID, &account.State)
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterGenesisAccounts(accountMgr *eth.AccountManager, accounts []models.GenesisAccount) error {
	log.Println("Registering genesis accounts")
	txs := make([]types.Transaction, 0, len(accounts))
	for i := range accounts {
		tx, err := accountMgr.RegisterAccount(&accounts[i].PublicKey)
		if err != nil {
			return errors.WithStack(err)
		}
		txs = append(txs, *tx)
	}

	receipts, err := chain.WaitForMultipleTxs(accountMgr.Blockchain.GetBackend(), txs...)
	if err != nil {
		return err
	}

	for i := range accounts {
		pubKeyID, err := accountMgr.RetrieveRegisteredPubKeyID(&receipts[i])
		if err != nil {
			return errors.WithStack(err)
		}

		if accounts[i].State.PubKeyID != *pubKeyID {
			return fmt.Errorf("different pubKeyID for account %s", accounts[i].PublicKey)
		}
	}

	return nil
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
