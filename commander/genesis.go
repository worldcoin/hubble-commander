package commander

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var (
	errGenesisAccountsUniqueStateID = fmt.Errorf("accounts must have unique state IDs")
)

func PopulateGenesisAccounts(storage *st.Storage, accounts []models.GenesisAccount) error {
	seenStateIDs := make(map[uint32]bool)
	for i := range accounts {
		account := &accounts[i]

		if seenStateIDs[account.StateID] {
			return errors.WithStack(errGenesisAccountsUniqueStateID)
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

	contractBatch, err := c.client.GetContractBatch(&batchID)
	if err != nil {
		return err
	}
	batch = contractBatch.ToModelBatch()
	batch.PrevStateRoot = *root

	return c.storage.AddBatch(batch)
}
