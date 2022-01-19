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

var (
	errGenesisAccountsUniqueStateID = fmt.Errorf("accounts must have unique state IDs")
	errMissingGenesisPublicKey      = fmt.Errorf("genesis accounts require public keys")
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

func RegisterGenesisAccountsAndCalculateTotalAmount(
	accountMgr *eth.AccountManager,
	accounts []models.GenesisAccount,
) (*models.Uint256, error) {
	log.Println("Registering genesis accounts")

	emptyPublicKey := models.PublicKey{}
	txs := make([]types.Transaction, 0, len(accounts))
	for i := range accounts {
		if accounts[i].PublicKey == emptyPublicKey {
			return nil, errors.WithStack(errMissingGenesisPublicKey)
		}

		tx, err := accountMgr.RegisterAccount(&accounts[i].PublicKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		txs = append(txs, *tx)
	}

	receipts, err := chain.WaitForMultipleTxs(accountMgr.Blockchain.GetBackend(), txs...)
	if err != nil {
		return nil, err
	}

	totalGenesisAmount := models.NewUint256(0)
	for i := range accounts {
		pubKeyID, err := accountMgr.RetrieveRegisteredPubKeyID(&receipts[i])
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if accounts[i].State.PubKeyID != *pubKeyID {
			return nil, fmt.Errorf("different pubKeyID for account %s", accounts[i].PublicKey)
		}
		totalGenesisAmount = totalGenesisAmount.Add(&accounts[i].State.Balance)
	}

	return totalGenesisAmount, nil
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
