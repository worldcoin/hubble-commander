package commander

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) manageRollupLoop(isProposer bool) {
	rollupLoopRunning := c.isRollupLoopActive()
	if isProposer && !rollupLoopRunning && c.batchCreationEnabled {
		log.Debugf("Commander is an active proposer, starting rollupLoop")
		c.startRollupLoop()
	} else if !isProposer && rollupLoopRunning {
		log.Debugf("Commander is no longer an active proposer, stoppping rollupLoop")
		c.stopRollupLoop()
	}
}

func (c *Commander) startRollupLoop() {
	if c.isRollupLoopActive() {
		return
	}

	ctx, cancel := context.WithCancel(c.workersContext)
	c.startWorker("Rollup Loop", func() error { return c.rollupLoop(ctx) })
	c.cancelRollupLoop = cancel
	c.setRollupLoopActive(true)
}

func (c *Commander) stopRollupLoop() {
	if !c.isRollupLoopActive() {
		return
	}
	if c.cancelRollupLoop != nil {
		c.cancelRollupLoop()
	}
	c.setRollupLoopActive(false)
}

func (c *Commander) rollupLoop(ctx context.Context) (err error) {
	ticker := time.NewTicker(c.cfg.Rollup.BatchLoopInterval)
	defer ticker.Stop()

	currentBatchType := batchtype.Transfer

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err = c.rollupLoopIteration(ctx, &currentBatchType)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Commander) rollupLoopIteration(ctx context.Context, currentBatchType *batchtype.BatchType) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	err = c.unsafeRollupLoopIteration(ctx, currentBatchType)
	if errors.Is(err, executor.ErrNotEnoughDeposits) {
		return c.unsafeRollupLoopIteration(ctx, currentBatchType)
	}
	return err
}

func (c *Commander) unsafeRollupLoopIteration(ctx context.Context, currentBatchType *batchtype.BatchType) (err error) {
	err = validateStateRoot(c.storage)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.txPool.UpdateMempool()
	if err != nil {
		return err
	}

	rollupCtx := executor.NewRollupLoopContext(c.storage, c.client, c.cfg.Rollup, c.metrics, c.txPool.Mempool(), ctx, *currentBatchType)
	defer rollupCtx.Rollback(&err)

	switchBatchType(currentBatchType)

	var (
		batch            *models.Batch
		commitmentsCount *int
	)
	duration, err := metrics.MeasureDuration(func() error {
		batch, commitmentsCount, err = rollupCtx.CreateAndSubmitBatch()
		return err
	})

	var rollupError *executor.RollupError
	if errors.As(err, &rollupError) {
		rollupCtx.Rollback(&err)
		return c.handleRollupError(rollupError, rollupCtx.GetErrorsToStore())
	}
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.BatchBuildAndSubmissionDuration, prometheus.Labels{
		"type": metrics.BatchTypeToMetricsBatchType(batch.Type),
	})

	logNewBatch(batch, *commitmentsCount, duration)

	err = rollupCtx.Commit()
	if err != nil {
		return err
	}
	return c.storage.SetTransactionErrors(rollupCtx.GetErrorsToStore()...)
}

func switchBatchType(batchType *batchtype.BatchType) {
	switch *batchType {
	case batchtype.Transfer:
		*batchType = batchtype.Create2Transfer
	case batchtype.Create2Transfer:
		*batchType = batchtype.MassMigration
	case batchtype.MassMigration:
		*batchType = batchtype.Deposit
	case batchtype.Deposit:
		*batchType = batchtype.Transfer
	case batchtype.Genesis:
		panic("Not supported")
	}
}

func (c *Commander) handleRollupError(err *executor.RollupError, errorsToStore []models.TxError) error {
	if err.IsLoggable {
		log.Warnf("%+v", err)
	}

	if errors.Is(err, executor.ErrNotEnoughDeposits) {
		return err
	}

	return c.storage.SetTransactionErrors(errorsToStore...)
}

func logNewBatch(batch *models.Batch, commitmentsCount int, duration *time.Duration) {
	log.Printf(
		"Submitted a %s batch with %d commitment(s) on chain in %s. Batch ID: %d. Transaction hash: %v",
		batch.Type.String(),
		commitmentsCount,
		duration,
		batch.ID.Uint64(),
		batch.TransactionHash,
	)
}

func logLatestCommitment(latestCommitment *models.CommitmentBase) {
	fields := log.Fields{
		"latestBatchID":      latestCommitment.ID.BatchID.String(),
		"latestCommitmentID": latestCommitment.ID.IndexInBatch,
	}
	log.WithFields(fields).Error("rollupLoop: Sanity check on state tree root failed")
}
