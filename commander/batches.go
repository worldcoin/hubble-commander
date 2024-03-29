package commander

import (
	"context"
	"errors"

	"github.com/Worldcoin/hubble-commander/commander/disputer"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
)

var ErrSyncedFraudulentBatch = errors.New("commander synced fraudulent batch")

func (c *Commander) syncBatches(ctx context.Context, startBlock, endBlock uint64) error {
	spanCtx, span := newBlockTracer.Start(ctx, "syncBatches")
	defer span.End()

	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	duration, err := metrics.MeasureDuration(func() error {
		return c.unsafeSyncBatches(spanCtx, startBlock, endBlock)
	})
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncBatchesMethod,
	})

	return nil
}

func (c *Commander) unsafeSyncBatches(ctx context.Context, startBlock, endBlock uint64) error {
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

	newRemoteBatches, err := c.client.GetBatches(
		ctx,
		&eth.BatchesFilters{
			StartBlockInclusive: startBlock,
			EndBlockInclusive:   &endBlock,
			FilterByBatchID:     filter,
		},
	)
	if err != nil {
		return err
	}

	for _, remoteBatch := range newRemoteBatches {
		err = c.syncRemoteBatch(ctx, remoteBatch)
		if err != nil {
			return err
		}

		err = c.syncPendingStakeWithdrawal(remoteBatch)
		if err != nil {
			return err
		}

		select {
		case <-c.workersContext.Done():
			return ErrIncompleteBlockRangeSync
		default:
		}
	}

	return nil
}

func (c *Commander) syncRemoteBatch(ctx context.Context, remoteBatch eth.DecodedBatch) error {
	var icError *syncer.InconsistentBatchError

	_, span := newBlockTracer.Start(ctx, "syncRemoteBatch")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("hubble.batchID", int64(remoteBatch.GetBase().ID.Uint64())),
		attribute.String("hubble.batchType", remoteBatch.GetBase().Type.String()),
		attribute.String("hubble.batchHash", remoteBatch.GetBase().Hash.String()),
		attribute.Int("hubble.commitmentCount", remoteBatch.GetCommitmentsLength()),
	)

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
		logFraudulentBatch(&remoteBatch.GetBase().ID, disputableErr.Reason)
		return c.disputeFraudulentBatch(remoteBatch.ToDecodedTxBatch(), disputableErr)
	}
	return err
}

func (c *Commander) syncBatch(remoteBatch eth.DecodedBatch) (err error) {
	syncCtx := syncer.NewContext(c.storage, c.client, c.cfg.Rollup, remoteBatch.GetBase().Type)
	defer syncCtx.Rollback(&err)

	err = syncCtx.SyncBatch(remoteBatch)
	if err != nil {
		return err
	}
	return syncCtx.Commit()
}

func (c *Commander) syncPendingStakeWithdrawal(remoteBatch eth.DecodedBatch) error {
	batchBase := remoteBatch.GetBase()
	if batchBase.Committer == c.blockchain.GetAccount().From {
		err := c.storage.AddPendingStakeWithdrawal(&models.PendingStakeWithdrawal{
			BatchID:           batchBase.ID,
			FinalisationBlock: batchBase.FinalisationBlock,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Commander) replaceBatch(localBatch *models.Batch, remoteBatch eth.DecodedBatch) error {
	log.Warnf("Reverting local batch(es) with ID(s) greater or equal to %s", localBatch.ID.String())

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
	executionCtx := executor.NewExecutionContext(c.storage, c.client, c.cfg.Rollup, c.metrics, context.Background())
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
