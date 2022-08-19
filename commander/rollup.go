package commander

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
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

	updateMempoolTicker := time.NewTicker(time.Second * 10)
	defer updateMempoolTicker.Stop()

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
		case <-updateMempoolTicker.C:
			// TODO: a long rollupLoopIteration can starve this metric,
			//       create a separate worker
			err := c.updateMempoolMetrics()
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
	spanCtx, span := rollupTracer.Start(ctx, "RollupLoop")
	defer span.End()

	err = validateStateRoot(c.storage)
	if err != nil {
		return err
	}

	rollupCtx := executor.NewRollupLoopContext(c.storage, c.client, c.cfg.Rollup, c.metrics, spanCtx, *currentBatchType)
	defer rollupCtx.Rollback(&err)
	span.SetAttributes(attribute.String("hubble.batchType", currentBatchType.String()))

	// this chooses the type of the next batch, currentBatchType is not read once
	// the rollupCtx has been created.
	switchBatchType(currentBatchType)

	var (
		batch            *models.Batch
		commitmentsCount *int
	)
	duration, err := metrics.MeasureDuration(func() error {
		batch, commitmentsCount, err = rollupCtx.CreateAndSubmitBatch(spanCtx)
		return err
	})
	if errors.Is(err, executor.ErrNotEnoughTxs) || errors.Is(err, executor.ErrNotEnoughDeposits) {
		// tell datadog to ignore this trace, we didn't do anything
		// this requires custom configuration of the dd agent:
		//  apm_config.filter_tags.reject = ["manual.drop:true"]
		// if we don't do this then ~ every 500µs we emit a new trace
		// https://docs.datadoghq.com/tracing/guide/ignoring_apm_resources/?tab=datadogyaml#ignoring-based-on-span-tags
		span.SetAttributes(attribute.Bool("manual.drop", true))
	}

	var rollupError *executor.RollupError
	if errors.As(err, &rollupError) {
		rollupCtx.Rollback(&err)
		return c.handleRollupError(rollupError)
	}
	if err != nil {
		return err
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.BatchBuildAndSubmissionDuration, prometheus.Labels{
		"type": metrics.BatchTypeToMetricsBatchType(batch.Type),
	})

	logNewBatch(batch, *commitmentsCount, duration)

	err = func() error {
		_, span := rollupTracer.Start(spanCtx, "rollupCtx.commit")
		defer span.End()

		return rollupCtx.Commit()
	}()
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) updateMempoolMetrics() error {
	allMempoolTxs, err := c.storage.GetAllMempoolTransactions()
	if err != nil {
		return err
	}

	transferCount, c2tCount, mmCount := 0, 0, 0
	for i := range allMempoolTxs {
		switch allMempoolTxs[i].TxType {
		case txtype.Transfer:
			transferCount += 1
		case txtype.Create2Transfer:
			c2tCount += 1
		case txtype.MassMigration:
			mmCount += 1
		default:
			panic("unknown tx type")
		}
	}

	c.metrics.MempoolSize.Set(float64(len(allMempoolTxs)))
	c.metrics.MempoolSizeTransfer.Set(float64(transferCount))
	c.metrics.MempoolSizeCreate2Transfer.Set(float64(c2tCount))
	c.metrics.MempoolSizeMassMigration.Set(float64(mmCount))

	return nil
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

func (c *Commander) handleRollupError(err *executor.RollupError) error {
	if err.IsLoggable {
		log.Warnf("%+v", err)
	}

	if errors.Is(err, executor.ErrNotEnoughDeposits) {
		return err
	}

	return nil
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
