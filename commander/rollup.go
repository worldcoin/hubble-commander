package commander

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrInvalidStateRoot = errors.New("current state tree root doesn't match latest commitment post state root")

func (c *Commander) manageRollupLoop(cancel context.CancelFunc, isProposer bool) context.CancelFunc {
	if isProposer && !c.rollupLoopRunning {
		log.Debugf("Commander is an active proposer, starting rollupLoop")
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
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

	currentBatchType := txtype.Transfer

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.stopChannel:
			return nil
		case <-ticker.C:
			err := c.rollupLoopIteration(ctx, &currentBatchType)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Commander) rollupLoopIteration(ctx context.Context, currentBatchType *txtype.TransactionType) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	err = validateStateRoot(c.storage)
	if err != nil {
		return errors.WithStack(err)
	}

	transactionExecutor, err := executor.NewTransactionExecutor(c.storage, c.client, c.cfg.Rollup, executor.TransactionExecutorOpts{Ctx: ctx})
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType, c.signaturesDomain)
	if *currentBatchType == txtype.Transfer {
		*currentBatchType = txtype.Create2Transfer
	} else {
		*currentBatchType = txtype.Transfer
	}
	if err != nil {
		var e *executor.RollupError
		if errors.As(err, &e) {
			return nil
		}
		return err
	}

	return transactionExecutor.Commit()
}

func validateStateRoot(storage *st.Storage) error {
	latestCommitment, err := storage.GetLatestCommitment()
	if st.IsNotFoundError(err) {
		return nil
	}
	if err != nil {
		return err
	}
	stateRoot, err := st.NewStateTree(storage).Root()
	if err != nil {
		return err
	}
	if latestCommitment.PostStateRoot != *stateRoot {
		log.WithFields(log.Fields{"latestBatchID": latestCommitment.IncludedInBatch.String()}).
			Debug("rollupLoop: Sanity check on state tree root in failed")
		return ErrInvalidStateRoot
	}
	return nil
}
