package commander

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrAccountLeavesInconsistency = fmt.Errorf("inconsistency in account leaves between the database and the contract")

// TODO extract event filtering logic to eth.Client

func (c *Commander) syncAccounts(start, end uint64) error {
	startTime := time.Now()

	newAccountsSingle, err := c.syncSingleAccounts(start, end)
	if err != nil {
		return err
	}
	newAccountsBatch, err := c.syncBatchAccounts(start, end)
	if err != nil {
		return err
	}

	newAccountsCount := *newAccountsSingle + *newAccountsBatch
	if newAccountsCount > 0 {
		measureAccountsSyncingDuration(startTime, c.metrics)
		logNewSyncedAccountsCount(newAccountsCount)
	}

	return nil
}

func (c *Commander) syncSingleAccounts(start, end uint64) (newAccountsCount *int, err error) {
	it, err := c.client.AccountRegistry.FilterSinglePubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	newAccountsCount = ref.Int(0)

	for it.Next() {
		tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistryABI.Methods["register"].ID) {
			continue // TODO handle internal transactions
		}

		account, err := c.client.ExtractSingleAccount(tx.Data(), it.Event)
		if err != nil {
			return nil, err
		}

		isNewAccount, err := saveSyncedSingleAccount(c.storage.AccountTree, account)
		if err != nil {
			return nil, err
		}
		if *isNewAccount {
			*newAccountsCount++
		}
	}
	return newAccountsCount, nil
}

func (c *Commander) syncBatchAccounts(start, end uint64) (newAccountsCount *int, err error) {
	it, err := c.client.AccountRegistry.FilterBatchPubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	newAccountsCount = ref.Int(0)

	for it.Next() {
		tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistryABI.Methods["registerBatch"].ID) {
			continue // TODO handle internal transactions
		}

		accounts, err := c.client.ExtractAccountsBatch(tx.Data(), it.Event)
		if err != nil {
			return nil, err
		}

		isNewAccount, err := saveSyncedBatchAccounts(c.storage.AccountTree, accounts)
		if err != nil {
			return nil, err
		}
		if *isNewAccount {
			*newAccountsCount += len(accounts)
		}
	}
	return newAccountsCount, nil
}

func saveSyncedSingleAccount(accountTree *storage.AccountTree, account *models.AccountLeaf) (isNewAccount *bool, err error) {
	err = accountTree.SetSingle(account)
	var accountExistsErr *storage.AccountAlreadyExistsError
	if errors.As(err, &accountExistsErr) {
		return ref.Bool(false), validateExistingAccounts(accountTree, *accountExistsErr.Account)
	}
	if err != nil {
		return nil, err
	}
	return ref.Bool(true), nil
}

func saveSyncedBatchAccounts(accountTree *storage.AccountTree, accounts []models.AccountLeaf) (isNewAccount *bool, err error) {
	err = accountTree.SetBatch(accounts)
	var accountBatchExistsErr *storage.AccountBatchAlreadyExistsError
	if errors.As(err, &accountBatchExistsErr) {
		return ref.Bool(false), validateExistingAccounts(accountTree, accountBatchExistsErr.Accounts...)
	}
	if err != nil {
		return nil, err
	}
	return ref.Bool(true), nil
}

func validateExistingAccounts(accountTree *storage.AccountTree, accounts ...models.AccountLeaf) error {
	for i := range accounts {
		existingAccount, err := accountTree.Leaf(accounts[i].PubKeyID)
		if err != nil {
			return err
		}
		if existingAccount.PublicKey != accounts[i].PublicKey {
			return errors.WithStack(ErrAccountLeavesInconsistency)
		}
	}
	return nil
}

func measureAccountsSyncingDuration(
	start time.Time,
	commanderMetrics *metrics.CommanderMetrics,
) time.Duration {
	duration := time.Since(start).Round(time.Millisecond)
	commanderMetrics.SyncingAccountsDuration.Observe(float64(duration.Milliseconds()))
	return duration
}

func logNewSyncedAccountsCount(newAccountsCount int) {
	log.Printf("Found %d new account(s)", newAccountsCount)
}
