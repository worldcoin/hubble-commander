package commander

import (
	"context"
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
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

	transactionExecutor, err := newTransactionExecutor(c.storage, c.client, c.cfg.Rollup, transactionExecutorOpts{ctx: ctx})
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType, c.signaturesDomain)
	if err != nil {
		var e *RollupError
		if errors.As(err, &e) {
			return nil
		}
		return err
	}

	if *currentBatchType == txtype.Transfer {
		*currentBatchType = txtype.Create2Transfer
	} else {
		*currentBatchType = txtype.Transfer
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

	err = t.submitBatch(batch, commitments)
	if err != nil {
		return err
	}

	log.Printf(
		"Submitted a %s batch with %d commitment(s) on chain in %s. Batch ID: %d. Transaction hash: %v",
		batchType.String(),
		len(commitments),
		time.Since(startTime).Round(time.Millisecond).String(),
		batch.ID.Uint64(),
		batch.TransactionHash,
	)
	return nil
}

func (t *transactionExecutor) buildTransferCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingTransfers(2*t.cfg.TxsPerCommitment, nil)
	if err != nil {
		return nil, err
	}
	return t.createTransferCommitments(pendingTransfers, domain)
}

func (t *transactionExecutor) buildCreate2TransfersCommitments(domain *bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := t.storage.GetPendingCreate2Transfers(2*t.cfg.TxsPerCommitment, nil)
	if err != nil {
		return nil, err
	}
	return t.createCreate2TransferCommitments(pendingTransfers, domain)
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

// TODO refactor to method on transactionExecutor?
func newPendingBatch(storage *st.Storage, batchType txtype.TransactionType) (*models.Batch, error) {
	stateTree := st.NewStateTree(storage)
	prevStateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}
	batchID, err := storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.Batch{
		ID:            *batchID,
		Type:          batchType,
		PrevStateRoot: prevStateRoot,
	}, nil
}
