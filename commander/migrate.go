package commander

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/client"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var errMissingBootstrapNodeURL = fmt.Errorf("bootstrap node URL is required for migration mode")

func (c *Commander) migrate() error {
	if c.cfg.Bootstrap.BootstrapNodeURL == nil {
		return errors.WithStack(errMissingBootstrapNodeURL)
	}

	hubbleClient := client.NewHubble(*c.cfg.Bootstrap.BootstrapNodeURL, c.cfg.API.AuthenticationKey)
	return c.migrateCommanderData(hubbleClient)
}

func (c *Commander) migrateCommanderData(hubble client.Hubble) error {
	//TODO: fetch pending txs

	err := c.syncFailedTxs(hubble)
	if err != nil {
		return err
	}

	err = c.syncPendingBatches(hubble)
	if err != nil {
		return err
	}

	c.setMigrate(false)
	return nil
}

func (c *Commander) syncFailedTxs(hubble client.Hubble) error {
	failedTxs, err := hubble.GetFailedTransactions()
	if err != nil {
		return err
	}

	if failedTxs.Len() > 0 {
		err = c.saveFailedTxs(failedTxs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commander) saveFailedTxs(failedTxs models.GenericTransactionArray) error {
	operations := make([]storage.DBOperation, failedTxs.Len())
	for i := 0; i < failedTxs.Len(); i++ {
		failedTx := failedTxs.At(i)
		operations[i] = func(txStorage *storage.TransactionStorage) error {
			return txStorage.AddTransaction(failedTx)
		}
	}

	dbTxsCount, err := c.storage.UpdateInMultipleTransactions(operations)
	if err != nil {
		err = fmt.Errorf("saving %d failed txs failed during database transaction #%d because of: %w", failedTxs.Len(), dbTxsCount, err)
		return errors.WithStack(err)
	}

	return nil
}

func (c *Commander) syncPendingBatches(hubble client.Hubble) error {
	pendingBatches, err := hubble.GetPendingBatches()
	if err != nil {
		return err
	}

	for i := range pendingBatches {
		err = c.syncPendingBatch(dtoToModelsBatch(&pendingBatches[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commander) syncPendingBatch(batch *models.PendingBatch) (err error) {
	ctx := executor.NewRollupLoopContext(c.storage, c.client, c.cfg.Rollup, c.metrics, context.Background(), batch.Type)
	defer ctx.Rollback(&err)

	err = ctx.ExecutePendingBatch(batch)
	if err != nil {
		return err
	}

	return ctx.Commit()
}

func dtoToModelsBatch(dtoBatch *dto.PendingBatch) *models.PendingBatch {
	batch := models.PendingBatch{
		ID:              dtoBatch.ID,
		Type:            dtoBatch.Type,
		TransactionHash: dtoBatch.TransactionHash,
		Commitments:     make([]models.PendingCommitment, 0, len(dtoBatch.Commitments)),
	}

	for i := range dtoBatch.Commitments {
		batch.Commitments = append(batch.Commitments, models.PendingCommitment(dtoBatch.Commitments[i]))
	}

	return &batch
}
