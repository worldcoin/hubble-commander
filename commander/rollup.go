package commander

import (
	"context"
	"errors"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
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
	transactionExecutor, err := newTransactionExecutorWithCtx(ctx, c.storage, c.client, c.cfg.Rollup)
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	if *currentBatchType == txtype.Transfer {
		err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType)
		*currentBatchType = txtype.Create2Transfer
	} else {
		err = transactionExecutor.CreateAndSubmitBatch(*currentBatchType)
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

func (t *transactionExecutor) CreateAndSubmitBatch(batchType txtype.TransactionType) error {
	var commitments []models.Commitment

	domain, err := t.storage.GetDomain(t.client.ChainState.ChainID)
	if err != nil {
		return err
	}

	if batchType == txtype.Transfer {
		commitments, err = buildTransferCommitments(t.storage, t.cfg, *domain)
	} else {
		commitments, err = buildCreate2TransfersCommitments(t.storage, t.client, t.cfg, *domain)
	}
	if err != nil {
		return err
	}

	err = submitBatch(t.ctx, batchType, commitments, t.storage, t.client, t.cfg)
	if err != nil {
		return err
	}
	return nil
}

func buildTransferCommitments(storage *st.Storage, cfg *config.RollupConfig, domain bls.Domain) ([]models.Commitment, error) {
	pendingTransfers, err := storage.GetPendingTransfers()
	if err != nil {
		return nil, err
	}
	return createTransferCommitments(pendingTransfers, storage, cfg, domain)
}

func buildCreate2TransfersCommitments(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	domain bls.Domain,
) ([]models.Commitment, error) {
	pendingTransfers, err := storage.GetPendingCreate2Transfers()
	if err != nil {
		return nil, err
	}
	return createCreate2TransferCommitments(pendingTransfers, storage, client, cfg, domain)
}
