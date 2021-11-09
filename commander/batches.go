package commander

import (
	"context"
	"errors"

	"github.com/Worldcoin/hubble-commander/commander/disputer"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
)

var ErrSyncedFraudulentBatch = errors.New("commander synced fraudulent batch")

func (c *Commander) syncBatches(startBlock, endBlock uint64) error {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	return c.unsafeSyncBatches(startBlock, endBlock)
}

func (c *Commander) unsafeSyncBatches(startBlock, endBlock uint64) error {
	latestBatchID, err := c.getLatestBatchID()
	if err != nil {
		return err
	}

	if c.invalidBatchID != nil && latestBatchID.Cmp(c.invalidBatchID) >= 0 {
		return ErrSyncedFraudulentBatch
	}

	filter := func(batchID *models.Uint256) bool {
		if batchID.Cmp(latestBatchID) <= 0 {
			log.Printf("Batch #%d already synced. Skipping...", batchID.Uint64())
			return false
		}
		if c.invalidBatchID != nil && batchID.Cmp(c.invalidBatchID) >= 0 {
			log.Printf("Batch #%d after dispute. Skipping...", batchID.Uint64())
			return false
		}
		return true
	}

	newRemoteBatches, err := c.client.GetBatches(&eth.BatchesFilters{
		StartBlockInclusive: startBlock,
		EndBlockInclusive:   &endBlock,
		FilterByBatchID:     filter,
	})
	if err != nil {
		return err
	}

	for i := range newRemoteBatches {
		err = c.syncRemoteBatch(newRemoteBatches[i].ToDecodedTxBatch())
		if err != nil {
			return err
		}

		select {
		case <-c.workersContext.Done():
			return ErrIncompleteBlockRangeSync
		default:
			continue
		}
	}

	return nil
}

func (c *Commander) syncRemoteBatch(remoteBatch eth.DecodedBatch) error {
	var icError *syncer.InconsistentBatchError

	err := c.syncOrDisputeRemoteBatch(remoteBatch)
	if errors.As(err, &icError) {
		return c.replaceBatch(icError.LocalBatch, remoteBatch)
	}
	return err
}

func (c *Commander) syncOrDisputeRemoteBatch(remoteBatch eth.DecodedBatch) error {
	var disputableErr *syncer.DisputableError

	err := c.syncBatch(remoteBatch)
	if errors.As(err, &disputableErr) {
		logFraudulentBatch(&remoteBatch.GetBatch().ID, disputableErr.Reason)
		return c.disputeFraudulentBatch(remoteBatch.ToDecodedTxBatch(), disputableErr)
	}
	return err
}

func (c *Commander) syncBatch(remoteBatch eth.DecodedBatch) (err error) {
	syncCtx := syncer.NewContext(c.storage, c.client, c.cfg.Rollup, remoteBatch.GetBatch().Type)
	defer syncCtx.Rollback(&err)

	err = syncCtx.SyncBatch(remoteBatch)
	if err != nil {
		return err
	}
	return syncCtx.Commit()
}

func (c *Commander) replaceBatch(localBatch *models.Batch, remoteBatch eth.DecodedBatch) error {
	log.WithFields(log.Fields{"batchID": localBatch.ID.String()}).
		Debug("Local batch inconsistent with remote batch, reverting local batch(es)")

	err := c.revertBatches(localBatch)
	if err != nil {
		return err
	}
	return c.syncOrDisputeRemoteBatch(remoteBatch)
}

func (c *Commander) disputeFraudulentBatch(
	remoteBatch *eth.DecodedTxBatch,
	disputableErr *syncer.DisputableError,
) (err error) {
	disputeCtx := disputer.NewContext(c.storage, c.client)

	switch disputableErr.Type {
	case syncer.Transition:
		err = disputeCtx.DisputeTransition(remoteBatch, disputableErr.CommitmentIndex, disputableErr.Proofs)
	case syncer.Signature:
		err = disputeCtx.DisputeSignature(remoteBatch, disputableErr.CommitmentIndex, disputableErr.Proofs)
	}
	if err != nil {
		return err
	}

	return ErrRollbackInProgress
}

func (c *Commander) revertBatches(startBatch *models.Batch) (err error) {
	executionCtx := executor.NewExecutionContext(c.storage, c.client, c.cfg.Rollup, context.Background())
	defer executionCtx.Rollback(&err)

	err = executionCtx.RevertBatches(startBatch)
	if err != nil {
		return err
	}
	return executionCtx.Commit()
}

func (c *Commander) getLatestBatchID() (*models.Uint256, error) {
	latestBatch, err := c.storage.GetLatestSubmittedBatch()
	if st.IsNotFoundError(err) {
		return models.NewUint256(0), nil
	} else if err != nil {
		return nil, err
	}
	return &latestBatch.ID, nil
}

func logFraudulentBatch(batchID *models.Uint256, reason string) {
	log.WithFields(log.Fields{"batchID": batchID.String()}).
		Infof("Found fraudulent batch. Reason: %s", reason)
}
