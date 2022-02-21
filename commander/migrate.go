package commander

import (
	"context"
	"fmt"

	"github.com/Worldcoin/hubble-commander/client"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var errMissingBootstrapNodeURL = fmt.Errorf("bootstrap node URL is required for migration mode")

func (c *Commander) migrate() error {
	nodeURL := c.cfg.Bootstrap.BootstrapNodeURL
	if nodeURL == nil {
		return errors.WithStack(errMissingBootstrapNodeURL)
	}
	log.Printf("Migration mode is on, syncing data from commander instance running at %s\n", *nodeURL)

	hubbleClient := client.NewHubble(*nodeURL, c.cfg.API.AuthenticationKey)
	return c.migrateCommanderData(hubbleClient)
}

func (c *Commander) migrateCommanderData(hubble client.Hubble) error {
	pendingTxsCount, err := c.syncPendingTxs(hubble)
	if err != nil {
		return err
	}
	log.Printf("Synced %d pending transaction(s)\n", pendingTxsCount)

	failedTxsCount, err := c.syncFailedTxs(hubble)
	if err != nil {
		return err
	}
	log.Printf("Synced %d failed transaction(s)\n", failedTxsCount)

	pendingBatchesCount, err := c.syncPendingBatches(hubble)
	if err != nil {
		return err
	}
	log.Printf("Synced %d pending batch(es)\n", pendingBatchesCount)

	c.setMigrate(false)
	return nil
}

func (c *Commander) syncPendingTxs(hubble client.Hubble) (int, error) {
	txs, err := hubble.GetPendingTransactions()
	if err != nil {
		return 0, err
	}

	err = c.storage.AddPendingTransactions(txs)
	if err != nil {
		return 0, err
	}

	return txs.Len(), nil
}

func (c *Commander) syncFailedTxs(hubble client.Hubble) (int, error) {
	txs, err := hubble.GetFailedTransactions()
	if err != nil {
		return 0, err
	}

	err = c.storage.AddFailedTransactions(txs)
	if err != nil {
		return 0, err
	}

	return txs.Len(), nil
}

func (c *Commander) syncPendingBatches(hubble client.Hubble) (int, error) {
	pendingBatches, err := hubble.GetPendingBatches()
	if err != nil {
		return 0, err
	}

	for i := range pendingBatches {
		err = c.syncPendingBatch(dtoToModelsBatch(&pendingBatches[i]))
		if err != nil {
			return 0, err
		}
	}

	return len(pendingBatches), nil
}

func (c *Commander) syncPendingBatch(batch *models.PendingBatch) (err error) {
	ctx := executor.NewRollupLoopContext(c.storage, c.client, c.cfg.Rollup, c.metrics, c.txPool.Mempool(), context.Background(), batch.Type)
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
		PrevStateRoot:   dtoBatch.PrevStateRoot,
		Commitments:     make([]models.PendingCommitment, 0, len(dtoBatch.Commitments)),
	}

	for i := range dtoBatch.Commitments {
		batch.Commitments = append(batch.Commitments, models.PendingCommitment(dtoBatch.Commitments[i]))
	}

	return &batch
}
