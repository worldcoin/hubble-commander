package commander

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var ErrInvalidStateRoot = errors.New("latest commitment state root doesn't match current one")

func (c *Commander) manageRollupLoop(cancel context.CancelFunc, isProposer bool) context.CancelFunc {
	if isProposer && !c.rollupLoopRunning {
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		c.startWorker(func() error { return c.rollupLoop(ctx) })
		c.rollupLoopRunning = true
	} else if !isProposer && c.rollupLoopRunning {
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

	transactionExecutor, err := NewTransactionExecutor(c.storage, c.client, c.cfg.Rollup, TransactionExecutorOpts{Ctx: ctx})
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	if *currentBatchType == txtype.Transfer {
		err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType, c.signaturesDomain)
		*currentBatchType = txtype.Create2Transfer
	} else {
		err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType, c.signaturesDomain)
		*currentBatchType = txtype.Transfer
	}
	if err != nil {
		var e *RollupError
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
		return ErrInvalidStateRoot
	}
	return nil
}
