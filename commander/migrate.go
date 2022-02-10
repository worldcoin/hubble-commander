package commander

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/client"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
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
	err := c.syncPendingTxs(hubble)
	if err != nil {
		return err
	}

	err = c.syncFailedTxs(hubble)
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

func (c *Commander) syncPendingTxs(hubble client.Hubble) error {
	txs, err := hubble.GetPendingTransactions()
	if err != nil {
		return err
	}

	return c.storage.AddPendingTransactions(txs)
}

func (c *Commander) syncFailedTxs(hubble client.Hubble) error {
	failedTxs, err := hubble.GetFailedTransactions()
	if err != nil {
		return err
	}

	if failedTxs.Len() > 0 {
		err = c.storage.SaveFailedTxs(failedTxs)
		if err != nil {
			return err
		}
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
