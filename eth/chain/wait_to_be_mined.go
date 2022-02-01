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
	pollInterval             = 500 * time.Millisecond
	ErrWaitToBeMinedTimedOut = fmt.Errorf("timeout on waiting for transaction to be mined")
)

func WaitToBeMined(r ReceiptProvider, timeout time.Duration, tx *types.Transaction) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r.Commit()
	return waitToBeMinedWithCtx(ctx, r, tx)
}

func waitToBeMinedWithCtx(ctx context.Context, r ReceiptProvider, tx *types.Transaction) (*types.Receipt, error) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		receipt, err := r.TransactionReceipt(ctx, tx.Hash())
		if err != nil && err != ethereum.NotFound {
			return nil, handleWaitToBeMinedError(err)
		}
		if receipt != nil && receipt.BlockNumber != nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, handleWaitToBeMinedError(ctx.Err())
		case <-ticker.C:
		}
	}
}

func handleWaitToBeMinedError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		err = errors.WithStack(ErrWaitToBeMinedTimedOut)
		log.Warnf("%+v", err)
		return err
	}
	return errors.WithStack(err)
}

type orderedReceipt struct {
	index   int
	receipt *types.Receipt
}

func WaitForMultipleTxs(r ReceiptProvider, timeout time.Duration, txs ...types.Transaction) ([]types.Receipt, error) {
	orChan := make(chan orderedReceipt, len(txs))
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r.Commit()

	group, ctx := errgroup.WithContext(ctxWithTimeout)
	for i := range txs {
		j := i
		group.Go(func() error {
			receipt, err := waitToBeMinedWithCtx(ctx, r, &txs[j])
			if err != nil {
				return err
			}
			orChan <- orderedReceipt{j, receipt}
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	close(orChan)
	result := make([]types.Receipt, len(txs))
	for or := range orChan {
		result[or.index] = *or.receipt
	}
	return result, nil
}

func CreateWaitForMultipleTxsHelper(r ReceiptProvider, timeout time.Duration) func(txs ...types.Transaction) ([]types.Receipt, error) {
	return func(txs ...types.Transaction) ([]types.Receipt, error) {
		return WaitForMultipleTxs(r, timeout, txs...)
	}
}
