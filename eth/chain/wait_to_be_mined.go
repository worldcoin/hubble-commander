package chain

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var (
	PollInterval             = 500 * time.Millisecond
	MineTimeout              = 5 * time.Minute
	ErrWaitToBeMinedTimedOut = fmt.Errorf("timeout on waiting for transaction to be mined")
)

func WaitToBeMined(r ReceiptProvider, tx *types.Transaction) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MineTimeout)
	defer cancel()

	return WaitToBeMinedCtx(ctx, r, tx)
}

func WaitToBeMinedCtx(ctx context.Context, r ReceiptProvider, tx *types.Transaction) (*types.Receipt, error) {
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		receipt, err := r.TransactionReceipt(ctx, tx.Hash())
		if err != nil && err != ethereum.NotFound {
			return nil, errors.WithStack(err)
		}
		if receipt != nil && receipt.BlockNumber != nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			err = ctx.Err()
			if errors.Is(err, context.DeadlineExceeded) {
				err = errors.WithStack(ErrWaitToBeMinedTimedOut)
				log.Warnf("%+v", err)
			}
			return nil, err
		case <-ticker.C:
		}
	}
}

func WaitForMultiple(r ReceiptProvider, txs []types.Transaction) ([]types.Receipt, error) {
	receiptChan := make(chan types.Receipt, len(txs))
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), MineTimeout)
	defer cancel()

	group, ctx := errgroup.WithContext(ctxWithTimeout)
	for i := range txs {
		j := i
		group.Go(func() error {
			receipt, err := WaitToBeMinedCtx(ctx, r, &txs[j])
			if err != nil {
				return err
			}
			receiptChan <- *receipt
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	close(receiptChan)
	receipts := make([]types.Receipt, 0, len(txs))
	for receipt := range receiptChan {
		receipts = append(receipts, receipt)
	}
	return receipts, nil
}
