package commander

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

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
	transactionExecutor, err := newTransactionExecutorWithCtx(ctx, c.storage, c.client, c.cfg.Rollup)
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

func (t *transactionExecutor) CreateAndSubmitBatch(batchType txtype.TransactionType, domain *bls.Domain) (err error) {
	startTime := time.Now()
	var commitments []models.Commitment
	batch, err := newPendingBatch(t.storage, batchType)
	if err != nil {
		return err
	}

	if batchType == txtype.Transfer {
		commitments, err = t.buildTransferCommitments(domain)
	} else {
		commitments, err = t.buildCreate2TransfersCommitments(domain)
	}
	if err != nil {
		return err
	}

	batch, err = t.submitBatch(batch, commitments)
	if err != nil {
		return err
	}

	log.Printf(
		"Submitted a %s batch with %d commitment(s) on chain in %s. Batch number: %d. Transaction hash: %v",
		batchType.String(),
		len(commitments),
		time.Since(startTime).Round(time.Millisecond).String(),
		batch.Number.Uint64(),
		batch.TransactionHash,
	)
	return nil
}

func (t *transactionExecutor) buildTransferCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingTransfers()
	if err != nil {
		return nil, err
	}
	return t.createTransferCommitments(pendingTransfers, domain)
}

func (t *transactionExecutor) buildCreate2TransfersCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingCreate2Transfers()
	if err != nil {
		return nil, err
	}
	return t.createCreate2TransferCommitments(pendingTransfers, domain)
}

func newPendingBatch(storage *st.Storage, batchType txtype.TransactionType) (*models.Batch, error) {
	stateTree := st.NewStateTree(storage)
	prevStateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}
	batchNumber, err := storage.GetNextBatchNumber()
	if err != nil {
		return nil, err
	}
	return &models.Batch{
		Type:              batchType,
		Number:            *batchNumber,
		PrevStateRootHash: prevStateRoot,
	}, nil
}
