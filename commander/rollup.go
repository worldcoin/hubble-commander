package commander

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *Commander) manageRollupLoop(cancel context.CancelFunc, isProposer bool) context.CancelFunc {
	if isProposer && !c.rollupLoopRunning {
		log.Debugf("Commander is an active proposer, starting rollupLoop")
		var ctx context.Context
		ctx, cancel = context.WithCancel(c.workersContext)
		c.startWorker(func() error { return c.rollupLoop(ctx) })
		c.rollupLoopRunning = true
	} else if !isProposer && c.rollupLoopRunning {
		log.Debugf("Commander is no longer an active proposer, stoppping rollupLoop")
		cancel()
		c.rollupLoopRunning = false
	}
	return cancel
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

	rollupCtx := executor.NewRollupLoopContext(c.storage, c.client, c.cfg.Rollup, ctx, *currentBatchType)
	defer rollupCtx.Rollback(&err)

	switchBatchType(currentBatchType)

	err = rollupCtx.CreateAndSubmitBatch()

	var rollupError *executor.RollupError
	if errors.As(err, &rollupError) {
		rollupCtx.Rollback(&err)
		return c.handleRollupError(rollupError, rollupCtx)
	}
	if err != nil {
		return err
	}

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
		*batchType = batchtype.Deposit
	case batchtype.Deposit:
		*batchType = batchtype.Transfer
	case batchtype.Genesis, batchtype.MassMigration:
		panic("Not supported")
	}
}

func (c *Commander) handleRollupError(rollupErr *executor.RollupError, rollupCtx executor.RollupLoopContext) error {
	if errors.Is(rollupErr, executor.ErrNotEnoughDeposits) {
		return rollupErr
	}

	if rollupErr.IsLoggable {
		log.Warnf("%+v", rollupErr)
	}

	return c.storage.SetTransactionErrors(rollupCtx.GetErrorsToStore()...)
}

func logLatestCommitment(latestCommitment *models.CommitmentBase) {
	fields := log.Fields{
		"latestBatchID":      latestCommitment.ID.BatchID.String(),
		"latestCommitmentID": latestCommitment.ID.IndexInBatch,
	}
	log.WithFields(fields).Error("rollupLoop: Sanity check on state tree root failed")
}
