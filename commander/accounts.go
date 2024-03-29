package commander

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var ErrAccountLeavesInconsistency = fmt.Errorf("inconsistency in account leaves between the database and the contract")

// TODO extract event filtering logic to eth.Client

func (c *Commander) syncAccounts(ctx context.Context, start, end uint64) error {
	var newAccountsSingle *int
	var newAccountsBatch *int

	spanCtx, span := newBlockTracer.Start(ctx, "syncAccounts")
	defer span.End()

	duration, err := metrics.MeasureDuration(func() (err error) {
		newAccountsSingle, err = c.syncSingleAccounts(start, end)
		if err != nil {
			return err
		}
		newAccountsBatch, err = c.syncBatchAccounts(spanCtx, start, end)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncAccountsMethod,
	})

	newAccountsCount := *newAccountsSingle + *newAccountsBatch
	logNewSyncedAccountsCount(newAccountsCount)

	return nil
}

func (c *Commander) syncSingleAccounts(start, end uint64) (newAccountsCount *int, err error) {
	it, err := c.getSinglePubKeyRegisteredIterator(start, end)
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

		if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistry.ABI.Methods["register"].ID) {
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

	return newAccountsCount, it.Error()
}

func (c *Commander) syncBatchAccounts(ctx context.Context, start, end uint64) (newAccountsCount *int, err error) {
	it, err := c.getBatchPubKeyRegisteredIterator(start, end)
	if err != nil {
		return nil, err
	}
	defer func() { _ = it.Close() }()

	newAccountsCount = ref.Int(0)

	for it.Next() {
		count, err := c.syncBatchAccountsTx(ctx, it.Event)
		if err != nil {
			return nil, err
		}
		*newAccountsCount += count
	}

	return newAccountsCount, it.Error()
}

func (c *Commander) syncBatchAccountsTx(
	ctx context.Context,
	event *accountregistry.AccountRegistryBatchPubkeyRegistered,
) (newAccountsCount int, err error) {
	_, span := newBlockTracer.Start(ctx, "syncBatchAccountsTx")
	defer span.End()

	tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), event.Raw.TxHash)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	span.SetAttributes(
		attribute.String("hubble.ethTx.hash", tx.Hash().String()),
		attribute.Int64("hubble.ethTx.nonce", int64(tx.Nonce())),
	)

	if !bytes.Equal(tx.Data()[:4], c.client.AccountRegistry.ABI.Methods["registerBatch"].ID) {
		// so we can alert if any of these appear
		span.SetAttributes(
			attribute.Bool("hubble.wasInternalTx", true),
		)
		span.SetStatus(codes.Error, "unhandled: internal account")
		return 0, nil // TODO handle internal transactions
	}

	accounts, err := c.client.ExtractAccountsBatch(tx.Data(), event)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	span.SetAttributes(
		attribute.Int("hubble.accountCount", len(accounts)),
	)

	isNewAccount, err := saveSyncedBatchAccounts(c.storage.AccountTree, accounts)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	span.SetStatus(codes.Ok, "")

	if *isNewAccount {
		return len(accounts), nil
	} else {
		return 0, nil
	}
}

func (c *Commander) getSinglePubKeyRegisteredIterator(start, end uint64) (*accountregistry.SinglePubKeyRegisteredIterator, error) {
	it := &accountregistry.SinglePubKeyRegisteredIterator{}

	err := c.client.FilterLogs(c.client.AccountRegistry.BoundContract, eth.SinglePubkeyRegisteredEvent, &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
}

func (c *Commander) getBatchPubKeyRegisteredIterator(start, end uint64) (*accountregistry.BatchPubKeyRegisteredIterator, error) {
	it := &accountregistry.BatchPubKeyRegisteredIterator{}

	err := c.client.FilterLogs(c.client.AccountRegistry.BoundContract, eth.BatchPubkeyRegisteredEvent, &bind.FilterOpts{
		Start: start,
		End:   &end,
	}, it)
	if err != nil {
		return nil, err
	}

	return it, nil
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

func logNewSyncedAccountsCount(newAccountsCount int) {
	if newAccountsCount > 0 {
		log.Printf("Found %d new account(s)", newAccountsCount)
	}
}
