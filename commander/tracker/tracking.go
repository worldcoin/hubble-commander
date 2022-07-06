package tracker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (t *Tracker) TrackSentTxs(ctx context.Context) error {
	wg := sync.WaitGroup{}
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg.Add(1)
	go func() {
		t.startReadingChannel(subCtx)
		wg.Done()
	}()

	errChan := make(chan error)
	wg.Add(1)
	go func() {
		err := t.startCheckingTxs(subCtx)
		errChan <- err
		wg.Done()
	}()

	err := <-errChan
	cancel()
	wg.Wait()
	return err
}

func (t *Tracker) startReadingChannel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case tx := <-t.txsChan:
			t.addTx(tx)
		}
	}
}

func (t *Tracker) startCheckingTxs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if t.isEmptyTxsQueue() {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			tx := t.firstTx()
			err := t.waitUntilTxMinedAndCheckForFail(tx)
			if err != nil {
				return err
			}

			t.removeFirstTx()
		}
	}
}

func (t *Tracker) waitUntilTxMinedAndCheckForFail(tx *types.Transaction) error {
	receipt, err := t.client.WaitToBeMined(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if receipt.Status == 1 {
		return nil
	}
	err = t.client.GetRevertMessage(tx, receipt)
	return fmt.Errorf("%w txHash=%s", err, tx.Hash().String())
}
